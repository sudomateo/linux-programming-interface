use std::env;
use std::error::Error;
use std::fmt::Display;
use std::fs::OpenOptions;
use std::io::{self, Read, Write};

#[derive(Debug)]
struct TeeError;

impl Error for TeeError {}

impl Display for TeeError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "did not write all bytes that were read")
    }
}

pub fn run() -> Result<(), Box<dyn Error>> {
    let mut args: Vec<String> = env::args().collect();

    let mut opts = OpenOptions::new();
    opts.create(true).write(true);
    if args.len() > 1 {
        if args.contains(&String::from("-a")) {
            if let Some(pos) = args.iter().position(|x| x == "-a") {
                args.remove(pos);
            };
            opts.append(true);
        } else {
            opts.truncate(true);
        }
    }

    let mut files = Vec::new();

    if args.len() > 1 {
        for i in 1..args.len() {
            let f = match opts.open(&args[i]) {
                Ok(f) => f,
                Err(e) => return Err(Box::new(e)),
            };
            files.push(f);
        }
    }

    let mut stdin = io::stdin();
    let mut stdout = io::stdout();
    let mut buf = [0; 1024];

    loop {
        let num_read = match stdin.read(&mut buf) {
            Ok(n) => n,
            Err(e) => return Err(Box::new(e)),
        };

        // No bytes were read.
        if num_read == 0 {
            break;
        }

        let num_written = match stdout.write(&buf[0..num_read]) {
            Ok(n) => n,
            Err(e) => return Err(Box::new(e)),
        };

        if num_read != num_written {
            return Err(Box::new(TeeError));
        }

        for mut file in files.iter() {
            let num_written = match file.write(&buf[0..num_read]) {
                Ok(n) => n,
                Err(e) => return Err(Box::new(e)),
            };

            if num_read != num_written {
                return Err(Box::new(TeeError));
            }
        }
    }

    Ok(())
}
