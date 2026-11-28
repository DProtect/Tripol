use futures_util::future::join_all;
use std::fs::OpenOptions;
use std::io::{self, BufWriter, Write};
use std::path::{Path, PathBuf};
use std::time::{Duration, Instant};

fn line() -> String {
    let mut line = String::new();
    io::stdin().read_line(&mut line).unwrap();
    line.trim().to_owned()
}

async fn receive(client: &reqwest::Client, url: &str, size: usize, destination: Option<&Path>) -> usize {
    let mut response = client.get(url).send().await.unwrap();
    assert_eq!(response.status(), 200);
    let length = match destination {
        Some(path) => {
            let output = OpenOptions::new().write(true).create_new(true).open(path).unwrap();
            let mut output = BufWriter::with_capacity(32768, output);
            let mut length = 0;
            while let Some(chunk) = response.chunk().await.unwrap() {
                output.write_all(&chunk).unwrap();
                length += chunk.len();
            }
            output.flush().unwrap();
            output.get_ref().sync_all().unwrap();
            length
        }
        None => response.text().await.unwrap().len(),
    };
    assert_eq!(length, size);
    length
}

fn main() {
    let runtime = tokio::runtime::Builder::new_current_thread().enable_all().build().unwrap();
    let client = reqwest::Client::builder().timeout(Duration::from_secs(10)).build().unwrap();
    let url = line();
    let count: usize = line().parse().unwrap();
    let mode = line();
    let width: usize = line().parse().unwrap();
    let size: usize = line().parse().unwrap();
    let directory = PathBuf::from(line());
    println!("READY");
    io::stdout().flush().unwrap();
    line();
    let mut samples = Vec::with_capacity(count / width);
    let start = Instant::now();
    let total: usize = runtime.block_on(async {
        let mut total = 0;
        if mode == "concurrent" {
            for _ in (0..count).step_by(width) {
                let began = Instant::now();
                let batch = join_all((0..width).map(|_| receive(&client, &url, size, None))).await;
                total += batch.into_iter().sum::<usize>();
                samples.push(began.elapsed().as_nanos());
            }
        } else {
            for index in 0..count {
                let path = (mode == "stream").then(|| directory.join(format!("{index}.bin")));
                let began = Instant::now();
                total += receive(&client, &url, size, path.as_deref()).await;
                samples.push(began.elapsed().as_nanos());
            }
        }
        total
    });
    println!("{} {}", start.elapsed().as_nanos(), total);
    for sample in samples {
        println!("{sample}");
    }
    println!("END");
    io::stdout().flush().unwrap();
    line();
}
