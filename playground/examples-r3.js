(function(root){const examples=[
  [
    "Hello world",
    "Hola mundo",
    "show \"Hello, TRIPOL!\"",
    "Hello, TRIPOL!\n",
    "",
    "Basics"
  ],
  [
    "Variables",
    "Variables",
    "name = \"TRIPOL\"\nversion = 27\nshow name, version",
    "TRIPOL 27\n",
    "",
    "Basics"
  ],
  [
    "Arithmetic",
    "Aritmética",
    "show 2 + 3 * 4\nshow (2 + 3) * 4\nshow -7 / 3, -7 % 3",
    "14\n20\n-2 -1\n",
    "",
    "Math"
  ],
  [
    "Constants",
    "Constantes",
    "const max = 100\nlet score:int = 42\nscore += 8\nshow score, max",
    "50 100\n",
    "",
    "Basics"
  ],
  [
    "Booleans",
    "Booleanos",
    "show true and not false\nshow false and (1 / 0 == 1)\nshow 7 >= 3",
    "true\nfalse\ntrue\n",
    "",
    "Basics"
  ],
  [
    "Conditionals",
    "Condicionales",
    "score = 85\nif score >= 90 {\n    show \"Excellent\"\n} else if score >= 60 {\n    show \"Passed\"\n} else {\n    show \"Try again\"\n}",
    "Passed\n",
    "",
    "Control"
  ],
  [
    "Range loop",
    "Bucle for",
    "for i in range(1, 6) {\n    show i\n}",
    "1\n2\n3\n4\n5\n",
    "",
    "Control"
  ],
  [
    "While loop",
    "Bucle while",
    "n = 3\nwhile n > 0 {\n    show n\n    n -= 1\n}\nshow \"Go!\"",
    "3\n2\n1\nGo!\n",
    "",
    "Control"
  ],
  [
    "Repeat",
    "Repetir",
    "total = 0\nrepeat 5 {\n    total += 2\n}\nshow total",
    "10\n",
    "",
    "Control"
  ],
  [
    "Functions",
    "Funciones",
    "func add(a:int, b:int)->int {\n    return a + b\n}\nshow add(12, 30)",
    "42\n",
    "",
    "Functions"
  ],
  [
    "Short function",
    "Función corta",
    "func square(n) => n * n\nshow square(9)",
    "81\n",
    "",
    "Functions"
  ],
  [
    "Factorial",
    "Factorial",
    "func factorial(n) {\n    if n <= 1 {return 1}\n    return n * factorial(n - 1)\n}\nshow factorial(6)",
    "720\n",
    "",
    "Algorithms"
  ],
  [
    "Fibonacci",
    "Fibonacci",
    "func fib(n) {\n    if n <= 1 {return n}\n    return fib(n - 1) + fib(n - 2)\n}\nshow fib(10)",
    "55\n",
    "",
    "Algorithms"
  ],
  [
    "Arrays",
    "Arreglos",
    "numbers = [10, 20, 30]\nshow numbers[1]\nshow len(numbers), sum(numbers)",
    "20\n3 60\n",
    "",
    "Collections"
  ],
  [
    "Dynamic arrays",
    "Arreglos dinámicos",
    "a = ints(5)\nfor i in range(5) {\n    a[i] = i * i\n}\nshow a",
    "[0, 1, 4, 9, 16]\n",
    "",
    "Collections"
  ],
  [
    "Array updates",
    "Modificar arreglos",
    "a = [1, 2, 3]\na[1] += 10\nshow a[1]",
    "12\n",
    "",
    "Collections"
  ],
  [
    "Loop control",
    "Control de bucles",
    "for i in range(10) {\n    if i == 2 {continue}\n    if i == 5 {break}\n    show i\n}",
    "0\n1\n3\n4\n",
    "",
    "Control"
  ],
  [
    "Console input",
    "Entrada de consola",
    "name = input(\"Your name: \")\nshow \"Hello\", name",
    "Your name: \nHello Ada\n",
    "Ada",
    "Basics"
  ],
  [
    "Assertions",
    "Aserciones",
    "value = 6 * 7\nassert value == 42, \"Unexpected value\"\nshow \"Test passed\"",
    "Test passed\n",
    "",
    "Basics"
  ],
  [
    "Prime checker",
    "Números primos",
    "func prime(n:int)->bool {\n    if n < 2 {return false}\n    for i in range(2, n) {\n        if n % i == 0 {return false}\n    }\n    return true\n}\nfor n in range(2, 12) {\n    if prime(n) {show n}\n}",
    "2\n3\n5\n7\n11\n",
    "",
    "Algorithms"
  ],
  [
    "String interpolation",
    "Interpolación",
    "name = \"TRIPOL\"\nscore = 20\nshow f\"Hello {name}! Score: {score * 2}\"",
    "Hello TRIPOL! Score: 40\n",
    "",
    "Text"
  ],
  [
    "Array iteration",
    "Recorrer arreglos",
    "numbers = [10, 20, 30]\nfor n in numbers {\n    show n\n}",
    "10\n20\n30\n",
    "",
    "Collections"
  ],
  [
    "Indexed iteration",
    "Recorrido con índices",
    "numbers = [7, 8, 9]\nfor i, value in numbers {\n    show f\"{i}: {value}\"\n}",
    "0: 7\n1: 8\n2: 9\n",
    "",
    "Collections"
  ],
  [
    "Even or odd",
    "Par o impar",
    "n = 17\nshow n % 2 == 0",
    "false\n",
    "",
    "Math"
  ],
  [
    "Absolute value",
    "Valor absoluto",
    "n = -42\nif n < 0 {n = -n}\nshow n",
    "42\n",
    "",
    "Math"
  ],
  [
    "Maximum of three",
    "Máximo de tres",
    "a = 7\nb = 20\nc = 11\nlargest = a\nif b > largest {largest = b}\nif c > largest {largest = c}\nshow largest",
    "20\n",
    "",
    "Math"
  ],
  [
    "Minimum of three",
    "Mínimo de tres",
    "a = 7\nb = 20\nc = -11\nsmallest = a\nif b < smallest {smallest = b}\nif c < smallest {smallest = c}\nshow smallest",
    "-11\n",
    "",
    "Math"
  ],
  [
    "Triangle area",
    "Área del triángulo",
    "base = 12\nheight = 7\nshow base * height / 2",
    "42\n",
    "",
    "Math"
  ],
  [
    "Rectangle area",
    "Área del rectángulo",
    "width = 8\nheight = 5\nshow width * height\nshow 2 * (width + height)",
    "40\n26\n",
    "",
    "Math"
  ],
  [
    "Integer power",
    "Potencia entera",
    "base = 3\nresult = 1\nrepeat 4 {result *= base}\nshow result",
    "81\n",
    "",
    "Math"
  ],
  [
    "Greatest common divisor",
    "Máximo común divisor",
    "a = 48\nb = 18\nwhile b != 0 {\n    r = a % b\n    a = b\n    b = r\n}\nshow a",
    "6\n",
    "",
    "Math"
  ],
  [
    "Least common multiple",
    "Mínimo común múltiplo",
    "func gcd(a:int,b:int)->int {\n    while b != 0 {\n        r = a % b\n        a = b\n        b = r\n    }\n    return a\n}\nshow 12 * 18 / gcd(12,18)",
    "36\n",
    "",
    "Math"
  ],
  [
    "Arithmetic progression",
    "Progresión aritmética",
    "for i in range(5) {show 3 + 4 * i}",
    "3\n7\n11\n15\n19\n",
    "",
    "Math"
  ],
  [
    "Geometric progression",
    "Progresión geométrica",
    "n = 2\nrepeat 5 {\n    show n\n    n *= 3\n}",
    "2\n6\n18\n54\n162\n",
    "",
    "Math"
  ],
  [
    "Integer square root",
    "Raíz cuadrada entera",
    "n = 80\nr = 0\nwhile (r+1)*(r+1) <= n {r += 1}\nshow r",
    "8\n",
    "",
    "Math"
  ],
  [
    "Digit sum",
    "Suma de dígitos",
    "n = 12345\ntotal = 0\nwhile n > 0 {\n    total += n % 10\n    n /= 10\n}\nshow total",
    "15\n",
    "",
    "Math"
  ],
  [
    "Binary digit count",
    "Cantidad de bits",
    "n = 37\nbits = 0\nwhile n > 0 {\n    bits += 1\n    n /= 2\n}\nshow bits",
    "6\n",
    "",
    "Math"
  ],
  [
    "Pythagorean triple",
    "Terna pitagórica",
    "a = 3\nb = 4\nc = 5\nshow a*a + b*b == c*c",
    "true\n",
    "",
    "Math"
  ],
  [
    "Sum from 1 to N",
    "Suma del 1 a N",
    "total = 0\nfor i in range(1,11) {total += i}\nshow total",
    "55\n",
    "",
    "Algorithms"
  ],
  [
    "Sum of squares",
    "Suma de cuadrados",
    "total = 0\nfor i in range(1,6) {total += i*i}\nshow total",
    "55\n",
    "",
    "Algorithms"
  ],
  [
    "Count divisors",
    "Contar divisores",
    "n = 24\ncount = 0\nfor i in range(1,n+1) {\n    if n%i == 0 {count += 1}\n}\nshow count",
    "8\n",
    "",
    "Algorithms"
  ],
  [
    "Perfect number",
    "Número perfecto",
    "n = 28\ntotal = 0\nfor i in range(1,n) {\n    if n%i == 0 {total += i}\n}\nshow total == n",
    "true\n",
    "",
    "Algorithms"
  ],
  [
    "Reverse digits",
    "Invertir dígitos",
    "n = 1234\nr = 0\nwhile n > 0 {\n    r = r*10+n%10\n    n /= 10\n}\nshow r",
    "4321\n",
    "",
    "Algorithms"
  ],
  [
    "Numeric palindrome",
    "Palíndromo numérico",
    "n = 1221\nx = n\nr = 0\nwhile x > 0 {\n    r = r*10+x%10\n    x /= 10\n}\nshow r == n",
    "true\n",
    "",
    "Algorithms"
  ],
  [
    "Collatz sequence",
    "Secuencia de Collatz",
    "n = 6\nwhile n != 1 {\n    show n\n    if n%2 == 0 {n /= 2} else {n = n*3+1}\n}\nshow n",
    "6\n3\n10\n5\n16\n8\n4\n2\n1\n",
    "",
    "Algorithms"
  ],
  [
    "Linear search",
    "Búsqueda lineal",
    "a = [7,3,9,2]\nfound = -1\nfor i,v in a {\n    if v == 9 {found = i; break}\n}\nshow found",
    "2\n",
    "",
    "Algorithms"
  ],
  [
    "Binary search",
    "Búsqueda binaria",
    "a = [2,4,6,8,10]\nlo = 0\nhi = len(a)-1\nfound = -1\nwhile lo <= hi {\n    mid = (lo+hi)/2\n    if a[mid] == 8 {found = mid; break}\n    if a[mid] < 8 {lo = mid+1} else {hi = mid-1}\n}\nshow found",
    "3\n",
    "",
    "Algorithms"
  ],
  [
    "Bubble sort",
    "Ordenamiento burbuja",
    "a = [5,1,4,2]\nfor i in range(len(a)) {\n    for j in range(len(a)-1-i) {\n        if a[j] > a[j+1] {\n            temp = a[j]\n            a[j] = a[j+1]\n            a[j+1] = temp\n        }\n    }\n}\nshow a",
    "[1, 2, 4, 5]\n",
    "",
    "Algorithms"
  ],
  [
    "Selection sort",
    "Ordenamiento selección",
    "a = [9,2,7,1]\nfor i in range(len(a)) {\n    smallest = i\n    for j in range(i+1,len(a)) {\n        if a[j] < a[smallest] {smallest = j}\n    }\n    temp = a[i]\n    a[i] = a[smallest]\n    a[smallest] = temp\n}\nshow a",
    "[1, 2, 7, 9]\n",
    "",
    "Algorithms"
  ],
  [
    "Insertion sort",
    "Ordenamiento inserción",
    "a = [4,3,1,2]\nfor i in range(1,len(a)) {\n    key = a[i]\n    j = i-1\n    while j >= 0 and a[j] > key {\n        a[j+1] = a[j]\n        j -= 1\n    }\n    a[j+1] = key\n}\nshow a",
    "[1, 2, 3, 4]\n",
    "",
    "Algorithms"
  ],
  [
    "Prefix sums",
    "Sumas acumuladas",
    "a = [3,1,4,2]\nb = ints(len(a))\ntotal = 0\nfor i,v in a {\n    total += v\n    b[i] = total\n}\nshow b",
    "[3, 4, 8, 10]\n",
    "",
    "Algorithms"
  ],
  [
    "Maximum subarray",
    "Subarreglo máximo",
    "a = [-2,3,-1,4,-5]\nbest = a[0]\ncurrent = 0\nfor v in a {\n    current += v\n    if current < v {current = v}\n    if current > best {best = current}\n}\nshow best",
    "6\n",
    "",
    "Algorithms"
  ],
  [
    "Run length count",
    "Conteo de rachas",
    "a = [1,1,2,2,2,3]\nruns = 1\nfor i in range(1,len(a)) {\n    if a[i] != a[i-1] {runs += 1}\n}\nshow runs",
    "3\n",
    "",
    "Algorithms"
  ],
  [
    "Array sum",
    "Suma de arreglo",
    "a = [4,8,12]\nshow sum(a)",
    "24\n",
    "",
    "Collections"
  ],
  [
    "Integer average",
    "Promedio entero",
    "a = [6,9,12]\nshow sum(a)/len(a)",
    "9\n",
    "",
    "Collections"
  ],
  [
    "Find largest item",
    "Mayor elemento",
    "a = [2,8,4]\nlargest = a[0]\nfor n in a {if n > largest {largest = n}}\nshow largest",
    "8\n",
    "",
    "Collections"
  ],
  [
    "Find smallest item",
    "Menor elemento",
    "a = [2,-8,4]\nsmallest = a[0]\nfor n in a {if n < smallest {smallest = n}}\nshow smallest",
    "-8\n",
    "",
    "Collections"
  ],
  [
    "Count occurrences",
    "Contar apariciones",
    "a = [1,2,1,3,1]\ncount = 0\nfor n in a {if n == 1 {count += 1}}\nshow count",
    "3\n",
    "",
    "Collections"
  ],
  [
    "Reverse array",
    "Invertir arreglo",
    "a = [1,2,3,4]\nfor i in range(len(a)/2) {\n    j = len(a)-1-i\n    temp = a[i]\n    a[i] = a[j]\n    a[j] = temp\n}\nshow a",
    "[4, 3, 2, 1]\n",
    "",
    "Collections"
  ],
  [
    "Array copy",
    "Copiar arreglo",
    "a = [3,5,7]\nb = ints(len(a))\nfor i,v in a {b[i] = v}\nshow b",
    "[3, 5, 7]\n",
    "",
    "Collections"
  ],
  [
    "Dot product",
    "Producto escalar",
    "a = [1,2,3]\nb = [4,5,6]\ntotal = 0\nfor i in range(len(a)) {total += a[i]*b[i]}\nshow total",
    "32\n",
    "",
    "Collections"
  ],
  [
    "Elementwise addition",
    "Suma de vectores",
    "a = [1,2,3]\nb = [4,5,6]\nc = ints(len(a))\nfor i in range(len(a)) {c[i] = a[i]+b[i]}\nshow c",
    "[5, 7, 9]\n",
    "",
    "Collections"
  ],
  [
    "Scale vector",
    "Escalar vector",
    "a = [2,4,6]\nfor i in range(len(a)) {a[i] *= 3}\nshow a",
    "[6, 12, 18]\n",
    "",
    "Collections"
  ],
  [
    "Histogram",
    "Histograma",
    "a = [0,1,1,2,2,2]\ncounts = ints(3)\nfor v in a {counts[v] += 1}\nshow counts",
    "[1, 2, 3]\n",
    "",
    "Collections"
  ],
  [
    "Matrix indexing",
    "Índices de matriz",
    "matrix = [1,2,3,4,5,6]\nrows = 2\ncols = 3\nfor r in range(rows) {\n    total = 0\n    for c in range(cols) {total += matrix[r*cols+c]}\n    show total\n}",
    "6\n15\n",
    "",
    "Collections"
  ],
  [
    "Diagonal sum",
    "Suma diagonal",
    "matrix = [1,2,3,4,5,6,7,8,9]\ntotal = 0\nfor i in range(3) {total += matrix[i*3+i]}\nshow total",
    "15\n",
    "",
    "Collections"
  ],
  [
    "Array rotation",
    "Rotar arreglo",
    "a = [1,2,3,4]\nfirst = a[0]\nfor i in range(len(a)-1) {a[i] = a[i+1]}\na[len(a)-1] = first\nshow a",
    "[2, 3, 4, 1]\n",
    "",
    "Collections"
  ],
  [
    "Count positive values",
    "Contar positivos",
    "a = [-1,2,0,4,-3]\ncount = 0\nfor v in a {if v > 0 {count += 1}}\nshow count",
    "2\n",
    "",
    "Collections"
  ],
  [
    "Join strings",
    "Unir textos",
    "first = \"Hello\"\nlast = \"TRIPOL\"\nshow first + \" \" + last",
    "Hello TRIPOL\n",
    "",
    "Text"
  ],
  [
    "String length",
    "Longitud de texto",
    "show len(\"TRIPOL\")",
    "6\n",
    "",
    "Text"
  ],
  [
    "Text equality",
    "Comparar textos",
    "show \"TRIPOL\" == \"TRIPOL\"\nshow \"TRIPOL\" != \"Rust\"",
    "true\ntrue\n",
    "",
    "Text"
  ],
  [
    "Escaped quotes",
    "Comillas escapadas",
    "show \"Hello \\\"world\\\"\"",
    "Hello \"world\"\n",
    "",
    "Text"
  ],
  [
    "Escaped braces",
    "Llaves escapadas",
    "show f\"{{score}} = {42}\"",
    "{score} = 42\n",
    "",
    "Text"
  ],
  [
    "Text banner",
    "Banner de texto",
    "line = \"\"\nrepeat 8 {line += \"*\"}\nshow line\nshow \"TRIPOL\"\nshow line",
    "********\nTRIPOL\n********\n",
    "",
    "Text"
  ],
  [
    "Numbered labels",
    "Etiquetas numeradas",
    "for i in range(1,4) {show f\"Item {i}\"}",
    "Item 1\nItem 2\nItem 3\n",
    "",
    "Text"
  ],
  [
    "Greeting function",
    "Función saludo",
    "func greet(name:string)->string {return \"Hello \" + name}\nshow greet(\"Ada\")",
    "Hello Ada\n",
    "",
    "Text"
  ],
  [
    "Two inputs",
    "Dos entradas",
    "name = input(\"Name: \")\ncity = input(\"City: \")\nshow f\"{name} from {city}\"",
    "Name: \nCity: \nAda from Bogota\n",
    "Ada\nBogota",
    "Text"
  ],
  [
    "Boolean interpolation",
    "Interpolación booleana",
    "active = true\nshow f\"Active: {active}\"",
    "Active: true\n",
    "",
    "Text"
  ],
  [
    "Temperature conversion",
    "Conversión de temperatura",
    "func fahrenheit(c:int)->int {return c*9/5+32}\nshow fahrenheit(20)",
    "68\n",
    "",
    "Functions"
  ],
  [
    "Function composition",
    "Composición de funciones",
    "func twice(n) => n*2\nfunc plusone(n) => n+1\nshow twice(plusone(5))",
    "12\n",
    "",
    "Functions"
  ],
  [
    "Boolean function",
    "Función booleana",
    "func even(n:int)->bool {return n%2==0}\nshow even(12)",
    "true\n",
    "",
    "Functions"
  ],
  [
    "Recursive countdown",
    "Cuenta regresiva recursiva",
    "func count(n:int)->int {\n    if n == 0 {return 0}\n    return 1+count(n-1)\n}\nshow count(8)",
    "8\n",
    "",
    "Functions"
  ],
  [
    "Recursive power",
    "Potencia recursiva",
    "func power(x:int,n:int)->int {\n    if n == 0 {return 1}\n    return x*power(x,n-1)\n}\nshow power(2,10)",
    "1024\n",
    "",
    "Functions"
  ],
  [
    "Return branches",
    "Retornos condicionales",
    "func sign(n:int)->int {\n    if n < 0 {return -1}\n    if n > 0 {return 1}\n    return 0\n}\nshow sign(-4), sign(0), sign(7)",
    "-1 0 1\n",
    "",
    "Functions"
  ],
  [
    "Local variables",
    "Variables locales",
    "func area(w:int,h:int)->int {\n    result = w*h\n    return result\n}\nshow area(6,7)",
    "42\n",
    "",
    "Functions"
  ],
  [
    "Forward function call",
    "Llamada adelantada",
    "show triple(7)\nfunc triple(n) => n*3",
    "21\n",
    "",
    "Functions"
  ],
  [
    "Nested calls",
    "Llamadas anidadas",
    "func add(a:int,b:int)->int {return a+b}\nshow add(add(1,2),add(3,4))",
    "10\n",
    "",
    "Functions"
  ],
  [
    "Function with array",
    "Función con arreglo",
    "func total(a:ints)->int {return sum(a)}\na = ints(3)\na[0]=1\na[1]=2\na[2]=3\nshow total(a)",
    "6\n",
    "",
    "Functions"
  ],
  [
    "Shopping receipt",
    "Recibo de compra",
    "prices = [1500,2000,3500]\nsubtotal = sum(prices)\ndiscount = subtotal*10/100\nshow f\"Subtotal: {subtotal}\"\nshow f\"Total: {subtotal-discount}\"",
    "Subtotal: 7000\nTotal: 6300\n",
    "",
    "Projects"
  ],
  [
    "Grade report",
    "Informe de notas",
    "grades = [80,90,70]\naverage = sum(grades)/len(grades)\nshow f\"Average: {average}\"\nif average >= 60 {show \"Passed\"} else {show \"Failed\"}",
    "Average: 80\nPassed\n",
    "",
    "Projects"
  ],
  [
    "Game score",
    "Puntuación de juego",
    "score = 0\npoints = [10,20,50]\nfor p in points {score += p}\nshow f\"Final score: {score}\"",
    "Final score: 80\n",
    "",
    "Projects"
  ],
  [
    "Health simulation",
    "Simulación de vida",
    "health = 100\ndamage = [15,30,10]\nfor d in damage {health -= d}\nshow f\"Health: {health}\"",
    "Health: 45\n",
    "",
    "Projects"
  ],
  [
    "Simple inventory",
    "Inventario simple",
    "stock = [10,5,8]\nstock[0] -= 2\nstock[1] += 3\nshow stock",
    "[8, 8, 8]\n",
    "",
    "Projects"
  ],
  [
    "Savings calculator",
    "Calculadora de ahorro",
    "balance = 0\nrepeat 12 {balance += 5000}\nshow f\"Savings: {balance}\"",
    "Savings: 60000\n",
    "",
    "Projects"
  ],
  [
    "Unit converter",
    "Conversor de unidades",
    "km = 7\nshow f\"{km} km = {km*1000} m\"",
    "7 km = 7000 m\n",
    "",
    "Projects"
  ],
  [
    "Trip cost",
    "Costo de viaje",
    "distance = 120\nfuelPerLiter = 15\nprice = 4000\nshow distance/fuelPerLiter*price",
    "32000\n",
    "",
    "Projects"
  ],
  [
    "ASCII rectangle",
    "Rectángulo ASCII",
    "for r in range(3) {show \"#####\"} ",
    "#####\n#####\n#####\n",
    "",
    "Projects"
  ],
  [
    "ASCII staircase",
    "Escalera ASCII",
    "line = \"\"\nrepeat 4 {\n    line += \"*\"\n    show line\n}",
    "*\n**\n***\n****\n",
    "",
    "Projects"
  ],
  [
    "Multiplication table",
    "Tabla de multiplicar",
    "n = 7\nfor i in range(1,6) {show f\"{n} x {i} = {n*i}\"}",
    "7 x 1 = 7\n7 x 2 = 14\n7 x 3 = 21\n7 x 4 = 28\n7 x 5 = 35\n",
    "",
    "Projects"
  ],
  [
    "Simple chatbot",
    "Chatbot simple",
    "question = input(\"You: \")\nif question == \"hello\" {show \"Bot: Hello!\"} else {show \"Bot: Try hello\"}",
    "You: \nBot: Hello!\n",
    "hello",
    "Projects"
  ]
];root.TripolExamples=examples;if(typeof module!=='undefined')module.exports=examples})(globalThis);
