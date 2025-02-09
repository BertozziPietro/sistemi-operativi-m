#include <mpi.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>

// Funzione per liberare tutta la memoria allocata
void free_all(float *m, float *r, int *send_counts, int *send_displs, int *recv_counts, int *recv_displs, float *input, float *output, int rank, bool all_malloc) {
	if (!rank) {
		free(m);
		free(r);
	}
	if (all_malloc) {
		free(send_counts);
		free(send_displs);
		free(recv_counts);
		free(recv_displs);
		free(input);
		free(output);
	}
}

int main(int argc, char **argv) {
	// Inizializzazione di MPI
	if (MPI_Init(&argc, &argv) != MPI_SUCCESS) {
		printf("Errore in MPI_Init\n");
		MPI_Abort(MPI_COMM_WORLD, 1);
	}
	
	// Dichiarazione delle variabili per la misurazione del tempo
	double start_computation = 0, start_total = 0, end_computation = 0, end_total = 0;
	int rank, size;
	MPI_Comm_rank(MPI_COMM_WORLD, &rank);
	MPI_Comm_size(MPI_COMM_WORLD, &size);
	
	// Inizio della misurazione del tempo totale nel processo 0
	if (rank == 0) start_total = MPI_Wtime();
	
	// Controllo degli argomenti
	if (argc != 2) {
		printf("Utilizzo: %s <dimensione>\n", argv[0]);
		return 1;
	}
	
	// Calcolo delle dimensioni della matrice
	int dim = atoi(argv[1]);
	int dim_res = !(dim % 2) ? dim / 2 : (dim + 1) / 2;
	int dim_row = !(dim % 2) ? dim + 1 : dim + 2;
	int len_r, len_m;
	
	// Verifica che la dimensione sia valida
	if (dim <= 1) {
		printf("La dimensione deve essere un numero intero maggiore di 1\n");
		return 1;
	}
	
	// Dichiarazione dei puntatori per la memoria allocata
	float *m = NULL, *r = NULL, *input = NULL, *output = NULL;
	int *send_counts = NULL, *send_displs = NULL, *recv_counts = NULL, *recv_displs = NULL;
	
	// Allocazione della matrice nel processo 0
	if (rank == 0) {
		len_m = dim_row * dim_row;
		len_r = dim_res * dim_res;
		m = (float *)malloc(len_m * sizeof(float));
		r = (float *)malloc(len_r * sizeof(float));
		if (m == NULL || r == NULL) {
			printf("Errore nell'allocazione di memoria del processo %d\n", rank);
			free_all(m, r, send_counts, send_displs, recv_counts, recv_displs, input, output, rank, false);
			MPI_Abort(MPI_COMM_WORLD, 1);
		}
		
		// Inizializzazione della matrice con valori casuali e bordi a zero
		for (int i = 0; i < len_m; i++) {
			if ((i < dim_row || !(i % (dim_row))) || ((dim % 2) && !((i + 1) % (dim + 2)))) m[i] = 0;
			else m[i] = (float)rand() / RAND_MAX;
		}
	}
	
	// Allocazione delle strutture per la distribuzione dei dati
	send_counts = (int *)malloc(size * sizeof(int));
	send_displs = (int *)malloc(size * sizeof(int));
	recv_counts = (int *)malloc(size * sizeof(int));
	recv_displs = (int *)malloc(size * sizeof(int));
	
	// Calcolo della suddivisione delle righe tra i processi
	int base_rows = dim_res / size;
	int extra_rows = dim_res % size;
	int max_rows = extra_rows != 0 ? base_rows + 1 : base_rows;
	
	// Allocazione della memoria per input e output
	input = malloc((max_rows * 2 + 1) * dim_row * sizeof(float));
	output = malloc(max_rows * dim_res * sizeof(float));
	if (send_counts == NULL || send_displs == NULL || recv_counts == NULL || recv_displs == NULL || input == NULL || output == NULL) {
		printf("Errore nell'allocazione di memoria\n");
		free_all(m, r, send_counts, send_displs, recv_counts, recv_displs, input, output, rank, true);
		MPI_Abort(MPI_COMM_WORLD, 1);
	}
	
	// Calcolo delle quantità di dati da inviare e ricevere
	int rows = base_rows + 1;
	bool extra = true;
	for (int k = 0, send_offset = 0, recv_offset = 0; k < size; k++) {
		if (extra && k == extra_rows) {
			rows--;
			extra = false;
		}
		send_counts[k] = (rows * 2 + 1) * (dim_row);
		send_displs[k] = send_offset;
		send_offset += (rows * 2) * (dim_row);
		recv_counts[k] = rows * dim_res;
		recv_displs[k] = recv_offset;
		recv_offset += rows * dim_res;
	}
	
	// Inizio della misurazione del tempo di calcolo nel processo 0
	MPI_Barrier(MPI_COMM_WORLD);
	if (rank == 0) start_computation = MPI_Wtime();
	
	// Distribuzione dei dati ai processi
	if (MPI_Scatterv(m, send_counts, send_displs, MPI_FLOAT, input, send_counts[rank], MPI_FLOAT, 0, MPI_COMM_WORLD) != MPI_SUCCESS) {
		printf("Errore in MPI_Scatterv\n");
		free_all(m, r, send_counts, send_displs, recv_counts, recv_displs, input, output, rank, true);
		MPI_Abort(MPI_COMM_WORLD, 1);
	}
	
	//if (rank == 0) start_computation = MPI_Wtime();
	
	// Computazione locale
	int input_rows = send_counts[rank] / (dim_row);
	for (int i = 1; i <= input_rows - 2; i += 2) {
		for (int j = 1; j < dim_row; j += 2) {
			int index = i * (dim_row) + j;
			float sum = input[index] + input[index - 1] + input[index + 1] +
						input[index - dim_row] + input[index - dim_row - 1] + input[index - dim_row + 1] +
						input[index + dim_row] + input[index + dim_row - 1] + input[index + dim_row + 1];
			output[j / 2 + ((i - 1) / 2 * dim_res)] = sum / 9.0;
		}
	}
	
	//if (rank == 0) end_computation = MPI_Wtime();
	
	// Raccolta dei risultati parziali
	if (MPI_Gatherv(output, recv_counts[rank], MPI_FLOAT, r, recv_counts, recv_displs, MPI_FLOAT, 0, MPI_COMM_WORLD) != MPI_SUCCESS) {
		printf("Errore in MPI_Gatherv\n");
		free_all(m, r, send_counts, send_displs, recv_counts, recv_displs, input, output, rank, true);
		MPI_Abort(MPI_COMM_WORLD, 1);
	}
	
	/*
	if (rank == 0) {
		for (int i = 0; i < len_m; i++) {
			if (!(i % dim_row)) printf("\n");
			printf("%f\t", m[i]);
		}
		printf("\n");
		for (int i = 0; i < len_r; i++) {
			if (!(i % dim_res)) printf("\n");
			printf("%f\t", r[i]);
		}
		printf("\n");
	}*/
	
	// Fine della misurazione del tempo di calcolo nel processo 0
	MPI_Barrier(MPI_COMM_WORLD);
	if (rank == 0) end_computation = MPI_Wtime();
	
	// Liberazione della memoria allocata
	free_all(m, r, send_counts, send_displs, recv_counts, recv_displs, input, output, rank, true);
	
	// Fine della misurazione del tempo totale nel processo 0
	MPI_Barrier(MPI_COMM_WORLD);
	if (rank == 0) {
		end_total = MPI_Wtime();
		printf("%d;%d;%f;%f\n", dim, size, end_total - start_total, end_computation - start_computation);
	}
	
	// Termine di MPI
	if (MPI_Finalize() != MPI_SUCCESS) {
		printf("Errore in MPI_Finalize\n");
		MPI_Abort(MPI_COMM_WORLD, 1);
	}
	return 0;
}
