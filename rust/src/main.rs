use std::fs::File;
use std::io::{BufRead, BufReader};
use std::thread;
use std::sync::{Arc, Mutex};
use clap::Parser;

extern crate num_cpus;

/// Buscador de texto em um arquivo
#[derive(Parser)]
#[clap(version = "1.0.0", author = "Henrique Delazeri <hwdelazeri@inf.ufrgs.br>")]
struct Opts {
    /// Arquivo de texto a ser lido
    file: String,
    /// Texto a ser procurado
    text: String,
    #[clap(short, long)]
    producers: Option<usize>,
    #[clap(short, long)]
    consumers: Option<usize>
}

fn main() {
    let opts: Opts = Opts::parse();

    let num_producers = match opts.producers {
        Some(v) => v,
        None => num_cpus::get() / 2
    };

    let num_consumers = match opts.consumers {
        Some(v) => v,
        None => num_cpus::get() / 2
    };

    let file_path = opts.file;
    let text_to_find = opts.text;

    println!("Procurando {} no arquivo {}", text_to_find, file_path);

    let file = match File::open(&file_path) {
        Err(why) => panic!("Não foi possível abrir o arquivo {}: {}", file_path, why),
        Ok(file) => file,
    };

    let reader = BufReader::new(file);
    let reader_mut = Arc::new(Mutex::new(reader));

    let rx = {
        let (s, r) = chan::sync(0);

        for _ in 0..num_producers {
            let reader_mutex_clone = Arc::clone(&reader_mut);
            let s = s.clone();

            thread::spawn(move || {
                loop {
                    let mut reader_lock = reader_mutex_clone.lock().unwrap();

                    let mut line = String::new();
                    match reader_lock.read_line(&mut line) {
                        Ok(bytes) => {
                            if bytes == 0 {
                                break;
                            }

                            s.send(line);
                        },
                        Err(_) => {
                            continue;
                        }
                    }
                }
            });
        }

        r
    };

    let wg = chan::WaitGroup::new();
    let mut consumers = Vec::with_capacity(num_consumers);

    for _ in 0..num_consumers {
        wg.add(1);
        
        let wg = wg.clone();
        let rx = rx.clone();
        let text_to_find = text_to_find.clone();

        consumers.push(thread::spawn(move || {
            let mut count = 0;

            for line in rx {
                if line.contains(&text_to_find) {
                    count = count + 1;
                }
            }

            wg.done();

            count
        }));
    }

    wg.wait();

    let mut total = 0;

    for handle in consumers {
        let count = handle.join().unwrap();
        total = total + count;
    }

    println!("{} foi encontrado {} vez(es) no arquivo {}", text_to_find, total, file_path);
}
