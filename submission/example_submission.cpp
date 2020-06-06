#include <iostream>
#include<stdio.h>
#include<stdlib.h>
using namespace std;

int add(int x, int y) {
    return x+y;
}

int main(int argc, char *argv[]) {
    int x, y;
    cin >> x;
    cin >> y;
    cout << add(x, y) << endl;
}
