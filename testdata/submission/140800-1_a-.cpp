#include <iostream>
using namespace std;

int main() {
    int n;

    cin >> n;

    for (int i = 0; i < n; ++i) {
        int x, y;

        cin >> x;
        cin >> y;

        cout << x*y << endl;
    }
}