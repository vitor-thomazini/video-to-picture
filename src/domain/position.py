class Position:
    def __init__(self, x, y) -> None:
        self.x = int(x)
        self.y = int(y)

    def asTuple(self) -> tuple:
        return (self.x, self.y)
    