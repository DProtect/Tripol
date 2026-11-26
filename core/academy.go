package core

import (
 "fmt"
 "strconv"
)

type academyLesson struct{ es,en,source,output,taskES,taskEN string }
var academyLessons=[]academyLesson{
 {"Primer programa","First program","show \"Hola, TRIPOL!\"\n","Hola, TRIPOL!","show muestra texto en la terminal. Cambia el saludo.","show prints text to the terminal. Change the greeting."},
 {"Variables","Variables","name = \"Alex\"\nscore = 10\nshow f\"{name}: {score * 2}\"\n","Alex: 20","Las variables infieren su tipo. Cambia score y vuelve a ejecutar.","Variables infer their type. Change score and run again."},
 {"Condiciones","Conditions","score = 80\nif score >= 60 {\n    show \"Passed\"\n} else {\n    show \"Try again\"\n}\n","Passed","Cambia score a 40 para probar la otra rama.","Change score to 40 to test the other branch."},
 {"Bucles","Loops","total = 0\nfor n in range(1, 6) {\n    total += n\n}\nshow total\n","15","range excluye el limite final. Prueba sumar del 1 al 10.","range excludes its end. Try summing 1 through 10."},
 {"Funciones","Functions","func twice(n) => n * 2\nshow twice(7)\n","14","Una funcion puede devolver un resultado. Prueba otro numero.","A function can return a result. Try another number."},
 {"Listas","Lists","numbers = list([10, 20, 30])\nnumbers.add(40)\nshow numbers\n","[10, 20, 30, 40]","Una lista puede crecer. Los indices comienzan en cero.","A list can grow. Indices start at zero."},
 {"Estructuras","Records","struct Player {\n    score: int = 0\n}\np = Player()\np.score = 5\nshow p.score\n","5","Una estructura agrupa datos relacionados.","A record groups related data."},
 {"Bibliotecas","Libraries","import tripol.math\nshow math.square(7)\n","49","tripol libs enumera bibliotecas; tripol docs math muestra sus funciones.","tripol libs lists libraries; tripol docs math shows their functions."},
 {"Errores","Errors","import tripol.errors\ntry {\n    show 10 / 0\n} catch code {\n    show errors.message(code)\n}\n","Integer overflow or division by zero","try/catch recupera errores comprobados. No todos los fallos son recuperables.","try/catch handles checked errors. Not all failures are recoverable."},
}
func academyCLI(args []string){
 if len(args)<1||len(args)>2{die("usage: tripol learn es|en [list|LESSON]")}
 spanish:=args[0]=="es"
 if !spanish&&args[0]!="en"{die("choose es or en")}
 if len(args)==1||args[1]=="list"{
  fmt.Println("TRIPOL Academy")
  for i,l:=range academyLessons{title:=l.en;if spanish{title=l.es};fmt.Printf("%d. %s\n",i+1,title)}
  fmt.Printf("tripol learn %s 1\n",args[0]);return
 }
 n,e:=strconv.Atoi(args[1]);if e!=nil||n<1||n>len(academyLessons){die("lesson must be between 1 and %d",len(academyLessons))}
 l:=academyLessons[n-1]
 title,task:=l.en,l.taskEN
 if spanish{title,task=l.es,l.taskES}
 fmt.Printf("TRIPOL Academy %d/%d - %s\n\n",n,len(academyLessons),title)
 fmt.Print(l.source)
 if spanish{fmt.Printf("\nResultado esperado:\n%s\n\nEjercicio: %s\n",l.output,task);fmt.Println("Copia el ejemplo en main.tri y ejecuta: tripol run")}else{fmt.Printf("\nExpected output:\n%s\n\nExercise: %s\n",l.output,task);fmt.Println("Copy the example into main.tri and run: tripol run")}
 if n<len(academyLessons){fmt.Printf("\ntripol learn %s %d\n",args[0],n+1)}
}
