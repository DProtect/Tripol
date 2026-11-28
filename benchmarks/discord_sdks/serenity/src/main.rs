use serenity::builder::{CreateActionRow, CreateButton, CreateEmbed};
use serenity::model::application::ButtonStyle;
use std::io::{self, Write};
use std::time::Instant;

fn classic() -> String {
    let embed = CreateEmbed::new()
        .title("TRIPOL")
        .description("Discord")
        .color(0x00A0FF)
        .field("Version", "0.28", false);
    let row = CreateActionRow::Buttons(vec![CreateButton::new("confirm")
        .label("Confirm")
        .style(ButtonStyle::Success)]);
    serde_json::to_string(&(embed, [row])).unwrap()
}

fn read_number() -> usize {
    let mut line = String::new();
    io::stdin().read_line(&mut line).unwrap();
    line.trim().parse().unwrap()
}

fn main() {
    println!("READY");
    io::stdout().flush().unwrap();
    let count = read_number();
    assert_eq!(read_number(), 0);
    let start = Instant::now();
    let sum: usize = (0..count)
        .map(|_| std::hint::black_box(classic()).len())
        .sum();
    println!("{} {}", start.elapsed().as_nanos(), sum);
    println!("{}", classic());
    io::stdout().flush().unwrap();
    read_number();
}
