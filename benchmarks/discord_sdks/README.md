# Serialización local de SDK de Discord / Local Discord SDK serialization

## Reproducción / Reproduction

Linux x86-64, AMD EPYC 7763 con tres CPU visibles, kernel 6.18.57 y glibc 2.41. TRIPOL 0.28 de esta rama compilado con Go 1.23.12 y perfil `turbo`. Discord.py 2.6.4 con CPython 3.12.14, Discord.js 14.25.1 con Node 24.22.3 y Serenity 0.12.5 con Rust 1.85.0 en modo release.

The recorded run used Linux x86-64 and an AMD EPYC 7763 with three visible CPUs. The pinned Python, npm and Cargo dependencies are included. Build the TRIPOL compiler from the checkout being evaluated, then run the commands below in this directory.

```sh
npm ci --ignore-scripts
python3 -m venv .venv
.venv/bin/python -m pip install -r requirements.txt
cargo build --manifest-path serenity/Cargo.toml --release --locked
tripol build serialize.tri --target linux --profile turbo -o tripol-serializer
python3 run.py --runs 7 --count 20000
```

Cada muestra crea y serializa 20.000 veces el mismo embed y una fila con un botón. El segundo caso usa un Container V2 con Text Display. La prueba verifica datos, colores, flags y estructura de cada salida después de medir. Se rotó el orden de los programas en las siete muestras y se publica la mediana. El arranque se mide hasta que el proceso puede recibir trabajo. Se lee RSS en reposo y el pico de RSS en `/proc`. El tiempo de CPU incluye el arranque. La tasa representa payloads serializados por segundo.

Each run builds and serializes the same embed and button row 20,000 times. The second workload is a V2 Container with Text Display. The harness validates the equivalent output fields, rotates process order and records the median of seven samples. Startup ends when the process accepts work. Idle and peak resident memory come from `/proc`. CPU time includes startup. Payload throughput measures serialization.

## Resultados / Results

| Carga / Workload | SDK | Arranque / Startup ms | RSS reposo / Idle KiB | RSS pico / Peak KiB | CPU ms | Serialización / Serialization ms |
| --- | --- | ---: | ---: | ---: | ---: | ---: |
| Embed + botón | TRIPOL | 0.341 | 52 | 60 | 70.540 | 71.640 |
| Embed + botón | Discord.js | 397.736 | 76,032 | 99,300 | 1,025.284 | 549.675 |
| Embed + botón | Discord.py | 199.594 | 36,556 | 36,728 | 638.676 | 417.022 |
| Embed + botón | Serenity | 1.319 | 1,116 | 1,124 | 13.396 | 12.247 |
| Components V2 | TRIPOL | 0.393 | 52 | 60 | 14.213 | 13.850 |
| Components V2 | Discord.js | 344.214 | 75,768 | 96,440 | 627.236 | 118.191 |
| Components V2 | Discord.py | 187.608 | 36,548 | 36,792 | 424.027 | 194.822 |

Serenity 0.12.5 aporta el caso común de embed y botón. Su API de construcción publicada carece del Container V2 tipado usado por las otras tres implementaciones. En TRIPOL se eliminaron validaciones UTF-8 redundantes del código que ya contaba caracteres con validación. Las dos cargas pasaron de 85.129 a 71.640 ms y de 19.021 a 13.850 ms. [Antes](results-before-utf8.json) · [Después](results.json)

The Serenity 0.12.5 builder contributes to the common embed and button case. Its released builder lacks the typed V2 Container used by the other three implementations. TRIPOL removed redundant UTF-8 checks from paths that already validated when counting characters. The two workloads went from 85.129 to 71.640 ms and from 19.021 to 13.850 ms. [Before](results-before-utf8.json) · [After](results.json)

El adaptador ELF de TRIPOL ocupa 38.082 B y el ejecutable de Serenity 552.368 B. La instalación de `node_modules` para la prueba ocupa 19.556.278 B y `site-packages` del entorno de pruebas de Python ocupa 32.537.691 B. Estos son componentes medidos individualmente. El tamaño desplegado de un bot incluye los lanzadores, intérpretes y demás dependencias de ejecución, y requiere una comparación separada.

The TRIPOL ELF adapter is 38,082 B and the Serenity executable is 552,368 B. The installed Node dependencies are 19,556,278 B and Python's test environment site packages are 32,537,691 B. These measurements cover individual components. Deployed bot size includes launchers, interpreters and other runtime dependencies and needs a separate test.

Las pruebas conectadas, la latencia de comandos, el rendimiento del Gateway y la estabilidad prolongada requieren otra carga. / Connected tests, command latency, Gateway throughput and long-run stability require a separate workload. [Validación conectada / Live validation](../../docs/discord-validation.md)
