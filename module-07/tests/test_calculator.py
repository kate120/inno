"""Простой калькулятор для демонстрации CI/CD."""

class Calculator:
    """Класс с базовыми математическими операциями."""
    
    @staticmethod
    def add(a: float, b: float) -> float:
        """Сложение."""
        return a + b
    
    @staticmethod
    def subtract(a: float, b: float) -> float:
        """Вычитание."""
        return a - b
    
    @staticmethod
    def multiply(a: float, b: float) -> float:
        """Умножение."""
        return a * b
    
    @staticmethod
    def divide(a: float, b: float) -> float:
        """Деление."""
        if b == 0:
            raise ValueError("Нельзя делить на ноль!")
        return a / b
    
    @staticmethod
    def power(a: float, b: float) -> float:
        """Возведение в степень."""
        return a ** b