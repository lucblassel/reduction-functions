# Reduction functions

When working with sequence nucleic acid data, some errors are expected.  
A common sequencing error is multiple insertions of the same nucleotide as the one just read.  
In order to reduce the impact of such errors on downstream tasks such as mapping or assembly, homopolymer compression is
used on the reads.

## Homopolymer compression

Homopolymer compression is a reduction function that compresses stretches of repeated nucleotides to a single one *(ie
AAAAAAAAA -> A)*.  
This can remove the insertion sequencing errors but it might also destroy some signal carried by true repetitions.

The goal of this project is to evaluate other reduction functions to see if they could be better than homopolymer
reduction. 