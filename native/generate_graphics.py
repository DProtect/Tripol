# Maintainer-only regeneration; no external toolchain used by installed TRIPOL.
import pathlib, subprocess, tempfile, re
root=pathlib.Path(__file__).resolve().parent
with tempfile.TemporaryDirectory() as tmp:
 o=pathlib.Path(tmp)/'graphics.o'; b=pathlib.Path(tmp)/'code.bin'
 subprocess.run(['as','--64','-o',str(o),str(root/'graphics.S')],check=True)
 rel=subprocess.check_output(['objdump','-r',str(o)],text=True)
 out=['// Code generated from native/graphics.S; DO NOT EDIT.','package core','var graphicsCode=map[string]nativeTemplate{']
 for s in ('gl_close','gl_frame','pad_read'):
  subprocess.run(['objcopy','-O','binary','-j','.text.'+s,str(o),str(b)],check=True)
  m=re.search(r'RELOCATION RECORDS FOR \[\.text\.'+s+r'\]:\n(.*?)(?=\n\n|\Z)',rel,re.S)
  fixes=[]
  for at,kind,label in re.findall(r'([0-9a-f]{16})\s+(R_X86_64_\w+)\s+(\S+)',m[1]):
   assert kind in ('R_X86_64_PC32','R_X86_64_PLT32') and label.endswith('-0x0000000000000004'),(kind,label)
   label=label.split('-')[0];label='iat:'+label[3:] if label.startswith('iat') else label
   fixes.append('{%d,"%s"}'%(int(at,16),label))
  out.append('"%s": {[]byte{%s}, []fix{%s}},'%(s,','.join('0x%02x'%x for x in b.read_bytes()),','.join(fixes)))
 out.append('}')
 (root.parent/'core'/'graphics_code.go').write_text('\n'.join(out)+'\n')
