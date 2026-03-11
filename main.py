from manga_ocr import MangaOcr
from pathlib import Path



if __name__ == "__main__":
    mocr = MangaOcr()
    input_dir = Path("img-teste/bubble-texts")
    output_dir = Path("out-text/manga-ocr")
    output_dir.mkdir(exist_ok=True)
    for file in input_dir.iterdir():
        if file.is_file():
            print(f"Processing: {file.name}")
            text = mocr(str(file))
            out_file = output_dir / (file.stem + ".txt")
            with open(out_file, "w", encoding="utf-8") as f:
                f.write(text)