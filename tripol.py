#!/usr/bin/env python3
"""TRIPOL v0.1 ORIGIN: statically typed source-to-native bootstrap."""
import argparse
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tempfile
from dataclasses import dataclass

VERSION = '0.1.0'
class Error(Exception): pass
@dataclass(frozen=True)
class Token:
    text: str
    line: int
    col: int
    kind: str

PAT = re.compile(r'(?P<space>[ \t\r]+)|(?P<comment>//[^\n]*|\#[^\n]*)|(?P<newline>\n)|(?P<string>"(?:[^"\\\n]|\\.)*")|(?P<int>[0-9]+)|(?P<id>[A-Za-z_][A-Za-z_0-9]*)|(?P<op>==|!=|<=|>=|&&|\|\||->|[+*/%<>=!(){},:;\-])|(?P<bad>.)')
RESERVED = set('func return show let const if else while true false int bool string void'.split())

def lex(src):
    tokens=[]; line=1; col=1
    for m in PAT.finditer(src):
        s=m.group(); k=m.lastgroup
        if k=='bad': raise Error(f'{line}:{col}: unexpected character {s!r}')
        if k not in ('space','comment'):
            tokens.append(Token(s,line,col,k))
        if k=='newline': line+=1; col=1
        else: col+=len(s)
    tokens.append(Token('<eof>',line,col,'eof'))
    return tokens

class Parser:
    def __init__(self, src): self.ts=lex(src); self.i=0
    @property
    def t(self): return self.ts[self.i]
    def fail(self,msg): raise Error(f'{self.t.line}:{self.t.col}: {msg}')
    def take(self): t=self.t; self.i+=1; return t
    def eat(self,s):
        if self.t.text==s: self.take(); return True
        return False
    def need(self,s):
        if not self.eat(s): self.fail(f'expected {s!r}, got {self.t.text!r}')
    def name(self):
        if self.t.kind!='id' or self.t.text in RESERVED: self.fail('expected identifier')
        return self.take().text
    def typ(self):
        if self.t.text not in ('int','bool','string'): self.fail('expected int, bool or string')
        return self.take().text
    def skip(self):
        while self.t.text in ('\n',';'): self.take()
    def end(self):
        if self.t.text not in ('\n',';','}','<eof>'): self.fail('expected newline or semicolon')
        self.skip()
    def program(self):
        fs=[]; ss=[]; self.skip()
        while self.t.kind!='eof':
            if self.eat('func'):
                name=self.name(); self.need('('); args=[]
                if not self.eat(')'):
                    while True:
                        n=self.name(); t=self.typ() if self.eat(':') else 'int'; args.append((n,t))
                        if self.eat(')'): break
                        self.need(',')
                ret=self.typ() if self.eat('->') else 'int'
                fs.append((name,args,ret,self.block()))
            else: ss.append(self.stmt())
            self.skip()
        return fs,ss
    def block(self):
        self.skip(); self.need('{'); self.skip(); out=[]
        while not self.eat('}'):
            if self.t.kind=='eof': self.fail('unclosed block')
            out.append(self.stmt()); self.skip()
        return out
    def stmt(self):
        line=self.t.line
        if self.eat('show'): n=('show',self.expr())
        elif self.eat('return'): n=('return',self.expr())
        elif self.eat('if'):
            cond=self.expr(); yes=self.block(); self.skip(); no=self.block() if self.eat('else') else []
            return (line,'if',cond,yes,no)
        elif self.eat('while'): return (line,'while',self.expr(),self.block())
        elif self.t.text in ('let','const'):
            mut=self.take().text=='let'; name=self.name(); ty=self.typ() if self.eat(':') else None
            self.need('='); n=('decl',name,ty,mut,self.expr())
        elif self.t.kind=='id' and self.ts[self.i+1].text=='=':
            name=self.name(); self.need('='); n=('assign',name,self.expr())
        else: n=('expr',self.expr())
        self.end(); return (line,*n)
    PREC={'||':1,'&&':2,'==':3,'!=':3,'<':4,'>':4,'<=':4,'>=':4,'+':5,'-':5,'*':6,'/':6,'%':6}
    def expr(self,minp=1):
        t=self.take(); s=t.text
        if s in ('-','!'): left=('un',s,self.expr(7))
        elif s=='(': left=self.expr(); self.need(')')
        elif t.kind=='int':
            v=int(s)
            if v>2**63-1: raise Error(f'{t.line}:{t.col}: integer outside signed 64-bit range')
            left=('int',v)
        elif t.kind=='string':
            try: v=json.loads(s); v.encode('utf-8')
            except (ValueError,UnicodeError): raise Error(f'{t.line}:{t.col}: invalid string') from None
            if '\0' in v: raise Error(f'{t.line}:{t.col}: NUL not supported in strings')
            left=('string',v)
        elif s in ('true','false'): left=('bool',s)
        elif t.kind=='id' and s not in RESERVED:
            if self.eat('('):
                args=[]
                if not self.eat(')'):
                    while True:
                        args.append(self.expr())
                        if self.eat(')'): break
                        self.need(',')
                left=('call',s,args)
            else: left=('var',s)
        else: raise Error(f'{t.line}:{t.col}: expected expression')
        while self.PREC.get(self.t.text,0)>=minp:
            op=self.take().text; right=self.expr(self.PREC[op]+1); left=('bin',op,left,right)
        return left

CT={'int':'int64_t','bool':'bool','string':'const char *'}
def cname(n): return 'tri_v_'+n
def fname(n): return 'tri_f_'+n
def cstr(s): return '"'+''.join('\\%03o'%b for b in s.encode('utf-8'))+'"'
RUNTIME=r'''#include <stdint.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>
static void tri_fail(void) { fputs("TRIPOL runtime error: invalid integer arithmetic\n", stderr); exit(1); }
static inline int64_t tri_add(int64_t a,int64_t b) { int64_t x; if(__builtin_add_overflow(a,b,&x)) tri_fail(); return x; }
static inline int64_t tri_sub(int64_t a,int64_t b) { int64_t x; if(__builtin_sub_overflow(a,b,&x)) tri_fail(); return x; }
static inline int64_t tri_mul(int64_t a,int64_t b) { int64_t x; if(__builtin_mul_overflow(a,b,&x)) tri_fail(); return x; }
static inline int64_t tri_div(int64_t a,int64_t b) { if(!b || (a==INT64_MIN && b==-1)) tri_fail(); return a/b; }
static inline int64_t tri_mod(int64_t a,int64_t b) { if(!b || (a==INT64_MIN && b==-1)) tri_fail(); return a%b; }
'''

class Emitter:
    def __init__(self,fs):
        self.sigs={}; self.scopes=[]; self.line=1; self.ret=None
        for n,args,t,_ in fs:
            if n in self.sigs: raise Error(f'duplicate function {n}')
            if len({a for a,_ in args})!=len(args): raise Error(f'duplicate parameter in {n}')
            self.sigs[n]=([ty for _,ty in args],t)
    def fail(self,msg): raise Error(f'line {self.line}: {msg}')
    def lookup(self,n):
        for s in reversed(self.scopes):
            if n in s: return s[n]
        return None
    def declare(self,n,t,mut):
        if n in self.scopes[-1]: self.fail(f'duplicate variable {n}')
        self.scopes[-1][n]=(t,mut)
    def require(self,a,b):
        if a!=b: self.fail(f'expected {b}, got {a}')
    def expr(self,e):
        k,*v=e
        if k=='int': return 'int',f'INT64_C({v[0]})'
        if k=='bool': return 'bool',v[0]
        if k=='string': return 'string',cstr(v[0])
        if k=='var':
            found=self.lookup(v[0])
            if found is None: self.fail(f'unknown variable {v[0]}')
            return found[0],cname(v[0])
        if k=='call':
            n,args=v
            if n not in self.sigs: self.fail(f'unknown function {n}')
            expected,ret=self.sigs[n]
            if len(args)!=len(expected): self.fail(f'{n} expects {len(expected)} arguments')
            cs=[]
            for a,t in zip(args,expected):
                actual,c=self.expr(a); self.require(actual,t); cs.append(c)
            return ret,f'{fname(n)}('+', '.join(cs)+')'
        if k=='un':
            op,e=v; t,c=self.expr(e); self.require(t,'bool' if op=='!' else 'int')
            return t,f'(!{c})' if op=='!' else f'tri_sub(0,{c})'
        op,l,r=v; a,x=self.expr(l); b,y=self.expr(r)
        if op in ('&&','||'):
            self.require(a,'bool'); self.require(b,'bool'); return 'bool',f'({x} {op} {y})'
        self.require(b,a)
        if op in ('==','!='):
            if a=='string': return 'bool',f'(strcmp({x},{y}) {op} 0)'
            return 'bool',f'({x} {op} {y})'
        self.require(a,'int')
        if op in ('<','>','<=','>='): return 'bool',f'({x} {op} {y})'
        return 'int',f"tri_{ {'+':'add','-':'sub','*':'mul','/':'div','%':'mod'}[op] }({x},{y})"
    def block(self,ss,indent=1,new=True):
        if new: self.scopes.append({})
        out=[]; pad='    '*indent
        for s in ss:
            self.line,k,*v=s
            if k in ('decl','assign'):
                if k=='decl': n,ty,mut,e=v; old=None
                else: n,e=v; old=self.lookup(n); ty=None; mut=True
                t,c=self.expr(e)
                if old:
                    if not old[1]: self.fail(f'cannot modify constant {n}')
                    self.require(t,old[0]); out.append(f'{pad}{cname(n)} = {c};')
                else:
                    if ty: self.require(t,ty)
                    # No shadowing: avoid C initializer binding to itself.
                    if self.lookup(n): self.fail(f'shadowing not supported: {n}')
                    self.declare(n,t,mut); out.append(f'{pad}{CT[t]} {cname(n)} = {c};')
            elif k=='show':
                t,c=self.expr(v[0]); call={'int':f'printf("%" PRId64 "\\n", {c})','bool':f'puts(({c}) ? "true" : "false")','string':f'puts({c})'}[t]; out.append(pad+call+';')
            elif k=='return':
                if self.ret is None: self.fail('return outside function')
                t,c=self.expr(v[0]); self.require(t,self.ret); out.append(f'{pad}return {c};')
            elif k=='expr':
                _,c=self.expr(v[0]); out.append(f'{pad}(void)({c});')
            elif k in ('if','while'):
                t,c=self.expr(v[0]); self.require(t,'bool'); out.append(f'{pad}{k} ({c}) {{'); out.extend(self.block(v[1],indent+1)); out.append(pad+'}')
                if k=='if' and v[2]: out[-1]+=' else {'; out.extend(self.block(v[2],indent+1)); out.append(pad+'}')
        if new: self.scopes.pop()
        return out
    @staticmethod
    def returns(ss):
        return any(s[1]=='return' or (s[1]=='if' and Emitter.returns(s[3]) and Emitter.returns(s[4])) for s in ss)
    def emit(self,fs,ss):
        lines=[RUNTIME]
        def signature(n,args,ret): return f'static {CT[ret]} {fname(n)}('+(', '.join(f'{CT[t]} {cname(a)}' for a,t in args) or 'void')+')'
        for n,args,ret,_ in fs: lines.append(signature(n,args,ret)+';')
        for n,args,ret,body in fs:
            self.scopes=[{a:(t,True) for a,t in args}]; self.ret=ret
            if not self.returns(body): raise Error(f'function {n} may finish without return')
            lines.append(signature(n,args,ret)+' {'); lines.extend(self.block(body,new=False)); lines.append('}')
        self.ret=None; self.scopes=[{}]; lines.append('int main(void) {'); lines.extend(self.block(ss,new=False)); lines.extend(['    return 0;','}']); return '\n'.join(lines)+'\n'

def translate(source):
    fs,ss=Parser(source).program(); return Emitter(fs).emit(fs,ss)

def compile_native(c,output,profile,cc):
    flags={'nano':['-Os','-ffunction-sections','-fdata-sections'],'balanced':['-O2'],'turbo':['-O3'],'debug':['-O0','-g']}[profile]
    if profile=='nano' and sys.platform!='darwin': flags+=['-Wl,--gc-sections','-s']
    with tempfile.TemporaryDirectory() as d:
        p=Path(d)/'program.c'; p.write_text(c,encoding='utf-8')
        subprocess.run([cc,'-std=c11',*flags,str(p),'-o',str(output)],check=True)

def main(argv=None):
    p=argparse.ArgumentParser(prog='tripol',description='TRIPOL ORIGIN native bootstrap'); p.add_argument('--version',action='version',version='TRIPOL '+VERSION)
    sub=p.add_subparsers(dest='cmd',required=True)
    for name in ('check','emit-c','build','run'):
        a=sub.add_parser(name); a.add_argument('source',type=Path)
        if name in ('emit-c','build'): a.add_argument('-o','--output',type=Path)
        if name in ('build','run'):
            a.add_argument('--profile',choices=('nano','balanced','turbo','debug'),default='nano'); a.add_argument('--cc',default=os.environ.get('CC') or shutil.which('clang') or shutil.which('gcc'))
    a=p.parse_args(argv)
    try:
        c=translate(a.source.read_text(encoding='utf-8'))
        if a.cmd=='check': print('OK'); return 0
        if a.cmd=='emit-c':
            if a.output: a.output.write_text(c,encoding='utf-8')
            else: print(c,end='')
            return 0
        if not a.cc: raise Error('GCC or Clang required to build native programs')
        suffix='.exe' if os.name=='nt' else ''
        if a.cmd=='build':
            out=(a.output or a.source.with_suffix(suffix)).resolve()
            if out==a.source.resolve(): raise Error('output must not overwrite source')
            compile_native(c,out,a.profile,a.cc); print(f'Built {out} ({out.stat().st_size} bytes)'); return 0
        with tempfile.TemporaryDirectory() as d:
            out=Path(d)/('program'+suffix); compile_native(c,out,a.profile,a.cc); return subprocess.run([str(out)]).returncode
    except (Error,OSError,UnicodeError,subprocess.CalledProcessError,RecursionError,ValueError) as e:
        print(f'tripol: {e}',file=sys.stderr); return 1

if __name__=='__main__': raise SystemExit(main())
