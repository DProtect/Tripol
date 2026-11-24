import sys
import tempfile
import subprocess
import unittest
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
import tripol

class TripolTests(unittest.TestCase):
    def execute(self, source, profile='nano'):
        with tempfile.TemporaryDirectory() as d:
            binary = Path(d)/('test.exe' if sys.platform=='win32' else 'test')
            tripol.compile_native(tripol.translate(source), binary, profile, 'gcc')
            return subprocess.run([str(binary)], capture_output=True, text=True)
    def output(self, source, expected):
        r=self.execute(source); self.assertEqual(r.returncode,0,r.stderr); self.assertEqual(r.stdout,expected)
    def invalid(self, source, message):
        with self.assertRaisesRegex(tripol.Error,message): tripol.translate(source)
    def test_hello(self): self.output('show "¡Hola!"','¡Hola!\n')
    def test_precedence(self): self.output('show 2 + 3 * 4\nshow (2 + 3) * 4','14\n20\n')
    def test_variables(self): self.output('x = 1\nx = x + 3\nshow x','4\n')
    def test_types(self): self.output('let x: bool = true\nconst s: string = "ok"\nshow x\nshow s','true\nok\n')
    def test_functions(self): self.output('func add(a,b) {return a+b}\nshow add(8,5)','13\n')
    def test_typed_functions(self): self.output('func name(x:string)->string {return x}\nshow name("yes")','yes\n')
    def test_recursion(self): self.output('func fib(n){if n<=1{return n}\nreturn fib(n-1)+fib(n-2)}\nshow fib(15)','610\n')
    def test_forward_call(self): self.output('show foo(5)\nfunc foo(x){return x*2}','10\n')
    def test_while(self): self.output('i=0\nwhile i<3 {show i; i=i+1}','0\n1\n2\n')
    def test_if_else(self): self.output('if false {show 1} else {show 2}','2\n')
    def test_bool(self): self.output('show !false && true\nshow false || true','true\ntrue\n')
    def test_short_circuit(self): self.output('show false && (1/0 == 1)','false\n')
    def test_strings(self): self.output('show "hola" == "hola"\nshow "a" != "b"','true\ntrue\n')
    def test_escape(self): self.output(r'show "a\nb\t\"c"','a\nb\t"c\n')
    def test_comments(self): self.output('# test\nshow 5 // test','5\n')
    def test_integer_division(self): self.output('show -7/3\nshow -7%3','-2\n-1\n')
    def test_return_branches(self): self.output('func f(x){if x==0{return 1}else{return 2}}\nshow f(0)','1\n')
    def test_scopes(self): self.output('x=0\nif true {let y=3; x=y}\nshow x','3\n')
    def test_const(self): self.invalid('const x=1\nx=2','constant')
    def test_wrong_type(self): self.invalid('let x: int = "a"','expected int')
    def test_unknown_var(self): self.invalid('show x','unknown variable')
    def test_unknown_func(self): self.invalid('show f(1)','unknown function')
    def test_wrong_args(self): self.invalid('func f(x){return x}\nshow f()','expects 1')
    def test_missing_return(self): self.invalid('func f(x){if x==0{return x}}','without return')
    def test_invalid_condition(self): self.invalid('if 1 {show 1}','expected bool')
    def test_shadowing(self): self.invalid('x=1\nif true {let x=2}','shadowing')
    def test_unclosed_block(self): self.invalid('if true {show 1','unclosed block')
    def test_bad_character(self): self.invalid('show @','unexpected character')
    def test_integer_literal_overflow(self): self.invalid('show 9223372036854775808','outside signed')
    def test_overflow(self):
        r=self.execute('show 9223372036854775807+1'); self.assertEqual(r.returncode,1); self.assertIn('runtime error',r.stderr)
    def test_div_zero(self): self.assertEqual(self.execute('show 1/0').returncode,1)
    def test_profiles(self):
        for p in ('balanced','turbo','debug'):
            with self.subTest(profile=p): self.assertEqual(self.execute('show 42',p).stdout,'42\n')
    def test_outside_scope(self): self.invalid('if true {x=1}\nshow x','unknown variable')
    def test_return_outside(self): self.invalid('return 1','outside function')

if __name__=='__main__': unittest.main()
