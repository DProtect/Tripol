# Primeros pasos con TRIPOL

Instala la distribución Windows x64 proporcionada por el mantenedor (sin firma digital), o compila desde el código fuente siguiendo [build.md](build.md). Abre un CMD nuevo. No necesitas instalar Go, Python o GCC para usar el compilador distribuido.

```cmd
tripol --version
tripol new MiProyecto
cd MiProyecto
tripol run
tripol build
main.exe
```

Edita `main.tri` en tu editor preferido. En CMD, `cd` cambia de carpeta y `cls` limpia la consola.

```tripol
nombre = input("Tu nombre?")
show "Hola", nombre
```

`tripol learn es` muestra el tutorial sin Internet. `tripol libs` lista bibliotecas; `tripol docs math` muestra su código. `tripol check main.tri --lang=es` comprueba errores. `tripol run main.tri --trace` añade una traza de funciones si el programa falla. Las nuevas funciones Windows (formulario y HTTPS) necesitan prueba real en Windows.

### Novedades de Infinity

```tripol
func saludo(nombre="Usuario", puntos=10)->string {
    return f"Hola {nombre}, tienes {puntos} puntos"
}
show saludo()
show saludo("Ana", 25)
```

Prueba `tripol libs`, `tripol docs collections` y `tripol ir main.tri`.

## Primer juego 2D (Windows x64, experimental)

```cmd
tripol new MiJuego2D --template game
cd MiJuego2D
tripol run
```

Usa WASD o las flechas y Esc para salir. El motor dibuja en un arreglo RGB y presenta la imagen mediante GDI. La ejecución de la ventana todavía requiere validación real en Windows. `tripol run examples/quantum_memory.tri` desde el código fuente comprueba la reserva y la reutilización de memoria. `tripol fix main.tri` muestra una sugerencia de cambio sin modificar código.
