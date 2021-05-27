# ClinCNV-Overlapper

Helps to calculate overlaps between samples. Works in two modes:

* **normal-only** would calculate overlaps within supplied samples list
* **normal and tumors mode** would calculate overlaps within each group of samples, and then it would get tumor results and filter any regions from normal one.

### Arguments

#### Required:
* `-n|--normals` or `-t|--tumors` expect a file with **tab**-separated values like `{SAMPLE}	normal/{SAMPLE}/{SAMPLE}_cnvs.tsv` and `SAMPLE	somatic/{TUMOR}-{NORMAL}/CNAs_{TUMOR}-{NORMAL}.txt` for germline and somatic ClinCNV results respectively for each sample, where `{SAMPLE}`, `{TUMOR}` and `{NORMAL}` are sample codes from coverage matrix. 
* `-o|--output` sets path for output file.

#### Optional:
* `-m|--min-overlap` sets filter for minimum overlap power to select final regions (default: 2).
* `--max-qvalue` provides filtration with maximum QValue (default: 1.0).
* `-b|--bed your_panel.bed` provides bed file for regions gene annotation.
* `--bed-gene-index` shows program what column is to be used from bed file for regions gene annotation (default: 4).
* `--buffer-size` sets buffer size in bytes for reading ClinCNV files (should be set and increased if necessary to avoid "take too long" error when reading ClinCNV files with long lines) (default: 100 * 1024).
* `-x|--genes-exclude-filter` provides multiline file with exclude patterns for gene names.
* `-f|--output-field-separator` provides a string for separation fields in genes and samples columns (default: " | ").
* `-s|--output-separator` provides a string for columns separation (default: "\t").

### Example

```shell
clincnv_overlapper \
    --normals data/normals_list.txt \
    [--tumors data/tumors_list.txt] \
    [--bed your_panel.bed] \
    [--bed-gene-index 5] \
    [--genes-exclude-filter data/gene_filters.txt] \
    --min-overlap 2 \
    [--buffer-size 1048576] \
    --max-qvalue 0.05 \
    --output-field-separator "; " \
    --output-separator "\t" \
    --output data/clincnv_overlaps.tsv
```

### TODO
Dependency `akamensky/argparse` would be soon replaced.