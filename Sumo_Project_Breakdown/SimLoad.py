from intersections import circle, stopStreet, tJunction, trafficLight


print("Welcome to Sumo!")


def showMenu():
    print("Select an instersection type:")
    print("1. Traffic circle")
    print("2. Stop street")
    print("3. T-Junction")
    print("4. Traffic Light")
    choice = input("Enter choice (1-4): ").strip()
    return choice


def getParams(tL: bool):
    trafficDensity = input("Enter traffic density (low/medium/high): ").strip().lower()
    if tL:
        greenTime = int(input("Enter green light duration in seconds: ").strip())
        redTime = int(input("Enter red light duration in seconds: ").strip())
        return {
            "Traffic Density": trafficDensity,
            "Green": greenTime,
            "Red": redTime
        }
    else:
        return {
            "Traffic Density": trafficDensity
        }


def main():
    tl = False
    choice = showMenu()

    if choice == '1':
        params = getParams(tl)
        circle.generate(params)
    elif choice == '2':
        params = getParams(tl)
        stopStreet.generate(params)
    elif choice == '3':
        params = getParams(tl)
        tJunction.generate(params)
    elif choice == '4':
        tl = True
        params = getParams(tl)
        trafficLight.generate(params)
    else:
        print("Invalid choice.")
        main()


if __name__ == "__main__":
    main()
    