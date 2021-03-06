from random import randint


def generate_constraints(width, height):
    f = open("constraints/generated_constraints.txt", "w")

    f.write(f"{width} {height} 30\n")

    for i in range(int((width * height) / 3)):
        f.write(f"parcel_{i} {randint(0, width - 1)} {randint(0, height - 1)} blue\n")

    f.write(f"palletTruck 0 {height - 1}\n")
    f.write(f"truck {width - 1} {randint(0, height - 1)} 4000 5\n")

    f.close()


if __name__ == '__main__':
    generate_constraints(20, 20)
