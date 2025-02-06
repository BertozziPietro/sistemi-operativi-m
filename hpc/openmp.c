#include <omp.h>
#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <stdbool.h>

float **create_matrix(int dim, bool oddFlag) {
    float **matrice = malloc(dim * sizeof(float *));
    for (int i = 0; i < dim; i++) matrice[i] = malloc(dim * sizeof(float));
    for (int i = 0; i < dim; i++) matrice[i][0] = matrice[0][i] = 0.0;
    for (int i = 1; i < dim; i++) for (int j = 1; j < dim; j++) {
    	if (oddFlag && j == (dim - 1)) matrice[i][j] = 0.0;
        else matrice[i][j] = (float)rand() / RAND_MAX;
    }
    return matrice;
}

void free_matrix(float **m, int dim) {
    for (int i = 0; i < dim; i++) free(m[i]);
    free(m);
}

void show_matrix(float **m, int dim) {
    for (int i = 0; i < dim; i++) {
        for (int j = 0; j < dim; j++) {
            float val = m[i][j];
            printf("%f\t", val);
        }
        printf("\n");
    }
    printf("\n");
}

float sum_adj(float **m, int i, int j) {
    float sum = 0;
    for (int ik = -1; ik < 2; ik++) for (int jk = -1; jk < 2; jk++) sum += m[i + ik][j + jk];
    return sum;
}

int main(int argc, char **argv) {
    double start_total = omp_get_wtime();
    if (argc < 3) {
        printf("Utilizzo: %s <dimensione>\n", argv[0]);
        return 1;
    }
    int dim = atoi(argv[1]);
    int dim_res = !(dim % 2) ? dim / 2 : (dim + 1) / 2;
    int dim_frm = !(dim % 2) ? dim + 1 : dim + 2;
    int size = atoi(argv[2]);
    if (dim <= 0) {
        printf("La dimensione deve essere un numero intero positivo.\n");
        return 1;
    }
    float **m = create_matrix(dim_frm, (dim % 2));
    float **r = malloc(dim * sizeof(float *));
    for (int i = 0; i < dim; i++)  r[i] = malloc(dim * sizeof(float));
    double start_computation = omp_get_wtime();
#pragma omp parallel for num_threads(size) shared(m, r, dim_res) schedule(static, 1)
	for (int i = 0; i < dim_res; i++) for (int j = 0; j < dim_res; j++) {
    	int i_m = i * 2 + 1;
        int j_m = j * 2 + 1;
        r[i][j] = sum_adj(m, i_m, j_m) / 9.0;
	}
    /*
    show_matrix(m, dim_frm);
    show_matrix(r, dim_res);	
    */
    double end_computation = omp_get_wtime();
    free_matrix(m, dim_frm);
    free_matrix(r, dim_res);
    double end_totale = omp_get_wtime();
    printf("%d;%d;%f;%f\n", dim, size, end_totale - start_total, end_computation - start_computation);
    return 0;
}
