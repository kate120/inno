// Main.java
public class Main {
    public static void main(String[] args) {
        System.out.println("Hello from Java in Docker!");
        System.out.println("Application started successfully!");
        
        // Бесконечный цикл, чтобы контейнер не завершался
        while (true) {
            try {
                Thread.sleep(5000);
                System.out.println("Still running...");
            } catch (InterruptedException e) {
                break;
            }
        }
    }
}