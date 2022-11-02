use std::process;

use tee;

fn main() {
    if let Err(e) = tee::run() {
        eprintln!("tee error: {e}");
        process::exit(1);
    }
}
