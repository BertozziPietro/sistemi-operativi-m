#!/bin/bash
#SBATCH --account=tra24_IngInfB2
#SBATCH --partition=g100_usr_prod
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=48
#SBATCH -o job.out
#SBATCH -e job.err
module load autoload intelmpi

for DIM in 2000 4000 8000 16000 32000 40000; do
	for I in {1..18}; do
		srun --ntasks=$I ./MPI $DIM
	done
	echo -e "\n" >> job.out
done
echo -e "\n" >> job.out
for RAP in 2000 2100 2200 2300 2400 2500; do
	for I in {1..18}; do
		SQRT_I=$(echo "scale=4; sqrt($I)" | bc -l)
		DIM=$(echo "$SQRT_I*$RAP" | bc)
		srun --ntasks=$I ./MPI $DIM
	done
	echo -e "\n" >> job.out
done
