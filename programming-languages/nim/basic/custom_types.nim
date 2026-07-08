type Color = enum
  Red, Green, Blue

var c: Color = Green
echo "Enum Color: ", c

type
  Point = object
    x, y: int

var p = Point(x: 10, y: 20)
echo "Object Point: ", p

type
  ShapeKind = enum
    Circle, Rectangle

  Shape = object
    case kind: ShapeKind
    of Circle:
      radius: float
    of Rectangle:
      width, height: float

var circ = Shape(kind: Circle, radius: 5.0)
var rect = Shape(kind: Rectangle, width: 10.0, height: 20.0)

echo "Circle radius: ", circ.radius
echo "Rectangle area: ", rect.width * rect.height
