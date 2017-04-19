int led=13;
int upto = 100;

void setup() {
  pinMode(12, OUTPUT);
  pinMode(11, OUTPUT);
  pinMode(10, OUTPUT);
  pinMode(9, OUTPUT);

  pinMode(13, OUTPUT);

  for (int i=9; i<13; ++i) {
    digitalWrite(i, LOW);
  }

  for (int i=0; i<10; ++i) {
    digitalWrite(12, HIGH);
    digitalWrite(12, LOW);
  }

  //for (int i=0; i<upto; ++i) {
  //  digitalWrite((i%4)+9, HIGH);
  //  digitalWrite((i%4)+9, LOW);
  //}
}

void loop() {
  digitalWrite(led, HIGH);
  delay(1000);
  digitalWrite(led, LOW);
  delay(1000);
}
