package core

// TIR SSA JOIN R1: a restricted terminal diamond has a single merged value.
// Push its final return through the two incoming edges, eliminating the
// phi's stack store/load without speculative evaluation. The original
// program has already been validated. Nonterminal or effectful diamonds
// continue through the ordinary checked backend.
func eliminateReturnPhi(body []stmt) []stmt {
 n:=len(body)
 if n<2{return body}
 branch,ret:=body[n-2],body[n-1]
 if branch.op!="if"||ret.op!="return"||ret.e==nil||ret.e.op!="var"||len(branch.yes)!=1||len(branch.no)!=1{return body}
 y,z:=branch.yes[0],branch.no[0]
 if y.op!="assign"||z.op!="assign"||y.name!=ret.e.s||z.name!=ret.e.s||y.typ!=""||z.typ!=""{return body}
 if !phiScalarExpr(y.e)||!phiScalarExpr(z.e){return body}
 // Each selected RHS executes at its original location and retains
 // checked overflow/division behavior. No expression is duplicated.
 yr,zr:=ret,ret
 yr.e,yr.line,yr.source=y.e,y.line,y.source
 zr.e,zr.line,zr.source=z.e,z.line,z.source
 branch.yes=[]stmt{yr};branch.no=[]stmt{zr}
 out:=append([]stmt(nil),body[:n-1]...)
 out[n-2]=branch
 return out
}

func phiScalarExpr(e *expr)bool{
 if e==nil{return false}
 switch e.op{
 case "int","bool","var":return true
 case "un","bin":
  for _,x:=range e.args{if !phiScalarExpr(x){return false}}
  return true
 }
 return false
}
