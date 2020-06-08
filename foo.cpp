//algo.cpp
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <iostream>
#include <math.h>

using namespace std;

extern "C" {
    char* cEater(int s) {

        size_t size = s * 1000000;

        char* out = new char[size];

        if (out == nullptr) {
            cout<<"Warning memory not allocated\n";
            return out;
        }

        for (int i=0; i <size; i++){
            *(out+i) = 'A';
        }
        

        return out;
    }
}
