use std::io::{self, Write};
use std::time::{Duration, Instant};

fn line() -> String {
    let mut line = String::new();
    io::stdin().read_line(&mut line).unwrap();
    line.trim().to_owned()
}

fn main() {
    let client = reqwest::blocking::Client::builder().timeout(Duration::from_secs(10)).build().unwrap();
    println!("READY");
    io::stdout().flush().unwrap();
    let url = line();
    let count: usize = line().parse().unwrap();
    let start = Instant::now();
    let mut total = 0;
    for _ in 0..count {
        let response = client.get(&url).send().unwrap();
        assert_eq!(response.status(), 200);
        let text = response.text().unwrap();
        assert_eq!(text.len(), 512);
        total += text.len();
    }
    println!("{} {}", start.elapsed().as_nanos(), total);
    io::stdout().flush().unwrap();
    line();
}
