#include <iostream>
using namespace std;

int main() {
  int N, K;
  string S;
  cin >> N >> K >> S;
  S[K - 1] = S[K - 1] + 'a' - 'A';
  cout << S << endl;
}
