(function(root){const examples=[
  [
    "Hello world",
    "Hola mundo",
    "show \"Hello, TRIPOL!\"",
    "Hello, TRIPOL!\n"
  ],
  [
    "Variables",
    "Variables",
    "name = \"TRIPOL\"\nversion = 27\nshow name, version",
    "TRIPOL 27\n"
  ],
  [
    "Arithmetic",
    "Aritmética",
    "show 2 + 3 * 4\nshow (2 + 3) * 4\nshow -7 / 3, -7 % 3",
    "14\n20\n-2 -1\n"
  ],
  [
    "Constants",
    "Constantes",
    "const max = 100\nlet score:int = 42\nscore += 8\nshow score, max",
    "50 100\n"
  ],
  [
    "Booleans",
    "Booleanos",
    "show true and not false\nshow false and (1 / 0 == 1)\nshow 7 >= 3",
    "true\nfalse\ntrue\n"
  ],
  [
    "Conditionals",
    "Condicionales",
    "score = 85\nif score >= 90 {\n    show \"Excellent\"\n} else if score >= 60 {\n    show \"Passed\"\n} else {\n    show \"Try again\"\n}",
    "Passed\n"
  ],
  [
    "Range loop",
    "Bucle for",
    "for i in range(1, 6) {\n    show i\n}",
    "1\n2\n3\n4\n5\n"
  ],
  [
    "While loop",
    "Bucle while",
    "n = 3\nwhile n > 0 {\n    show n\n    n -= 1\n}\nshow \"Go!\"",
    "3\n2\n1\nGo!\n"
  ],
  [
    "Repeat",
    "Repetir",
    "total = 0\nrepeat 5 {\n    total += 2\n}\nshow total",
    "10\n"
  ],
  [
    "Functions",
    "Funciones",
    "func add(a:int, b:int)->int {\n    return a + b\n}\nshow add(12, 30)",
    "42\n"
  ],
  [
    "Short function",
    "Función corta",
    "func square(n) => n * n\nshow square(9)",
    "81\n"
  ],
  [
    "Factorial",
    "Factorial",
    "func factorial(n) {\n    if n <= 1 {return 1}\n    return n * factorial(n - 1)\n}\nshow factorial(6)",
    "720\n"
  ],
  [
    "Fibonacci",
    "Fibonacci",
    "func fib(n) {\n    if n <= 1 {return n}\n    return fib(n - 1) + fib(n - 2)\n}\nshow fib(10)",
    "55\n"
  ],
  [
    "Arrays",
    "Arreglos",
    "numbers = [10, 20, 30]\nshow numbers[1]\nshow len(numbers), sum(numbers)",
    "20\n3 60\n"
  ],
  [
    "Dynamic arrays",
    "Arreglos dinámicos",
    "a = ints(5)\nfor i in range(5) {\n    a[i] = i * i\n}\nshow a",
    "[0, 1, 4, 9, 16]\n"
  ],
  [
    "Array updates",
    "Modificar arreglos",
    "a = [1, 2, 3]\na[1] += 10\nshow a[1]",
    "12\n"
  ],
  [
    "Loop control",
    "Control de bucles",
    "for i in range(10) {\n    if i == 2 {continue}\n    if i == 5 {break}\n    show i\n}",
    "0\n1\n3\n4\n"
  ],
  [
    "Console input",
    "Entrada de consola",
    "name = input(\"Your name: \")\nshow \"Hello\", name",
    "Your name: Ada\nHello Ada\n",
    "Ada"
  ],
  [
    "Assertions",
    "Aserciones",
    "value = 6 * 7\nassert value == 42, \"Unexpected value\"\nshow \"Test passed\"",
    "Test passed\n"
  ],
  [
    "Prime checker",
    "Números primos",
    "func prime(n:int)->bool {\n    if n < 2 {return false}\n    for i in range(2, n) {\n        if n % i == 0 {return false}\n    }\n    return true\n}\nfor n in range(2, 12) {\n    if prime(n) {show n}\n}",
    "2\n3\n5\n7\n11\n"
  ],
  [
    "String interpolation",
    "Interpolación",
    "name = \"TRIPOL\"\nscore = 20\nshow f\"Hello {name}! Score: {score * 2}\"",
    "Hello TRIPOL! Score: 40\n"
  ],
  [
    "Array iteration",
    "Recorrer arreglos",
    "numbers = [10, 20, 30]\nfor n in numbers {\n    show n\n}",
    "10\n20\n30\n"
  ],
  [
    "Indexed iteration",
    "Recorrido con índices",
    "numbers = [7, 8, 9]\nfor i, value in numbers {\n    show f\"{i}: {value}\"\n}",
    "0: 7\n1: 8\n2: 9\n"
  ]
];root.TripolExamples=examples;if(typeof module!=='undefined')module.exports=examples})(globalThis);
