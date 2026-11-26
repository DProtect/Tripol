package core
import(
 "fmt"
 "testing"
)
func TestAcademyExamples(t *testing.T){
 for i,l:=range academyLessons{
  t.Run(fmt.Sprint(i+1),func(t *testing.T){runCase(t,l.source,l.output+"\n",false)})
 }
}
