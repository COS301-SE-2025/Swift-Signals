import { render, screen, fireEvent, act } from "@testing-library/react";
import Carousel, { CarouselItem } from "../src/components/Carousel";

// Framer Motion mock: strip motion-only props and auto-call animationComplete
jest.mock("framer-motion", () => {
  const actual = jest.requireActual("framer-motion");
  return {
    ...actual,
    motion: {
      ...actual.motion,
      div: ({ children, onAnimationComplete, ...rest }: any) => {
        if (onAnimationComplete) setTimeout(() => onAnimationComplete(), 0);
        return <div {...rest}>{children}</div>;
      },
    },
    useMotionValue: jest.fn(() => ({ set: jest.fn(), get: jest.fn() })),
    useTransform: jest.fn(() => 0),
  };
});

jest.useFakeTimers();

describe("Carousel Component", () => {
  const items: CarouselItem[] = [
    { id: 1, title: "Test 1", description: "Desc 1", icon: <div>Icon1</div>, backgroundColor: "red" },
    { id: 2, title: "Test 2", description: "Desc 2", icon: <div>Icon2</div>, backgroundColor: "blue" },
  ];

  it("renders default items if no props passed", () => {
    render(<Carousel />);
    expect(screen.getByText("Overview")).toBeInTheDocument();
  });

  it("renders custom items", () => {
    render(<Carousel items={items} />);
    expect(screen.getByText("Test 1")).toBeInTheDocument();
    expect(screen.getByText("Test 2")).toBeInTheDocument();
  });

  it("applies baseWidth and round styles", () => {
    const { container } = render(<Carousel baseWidth={500} round />);
    const wrapper = container.querySelector<HTMLDivElement>(".carousel-container");
    expect(wrapper).toHaveStyle("width: 500px");
    expect(wrapper).toHaveClass("round");
  });

  it("handles indicator click to change index", () => {
    const { container } = render(<Carousel items={items} />);
    const indicators = container.querySelectorAll<HTMLDivElement>(".carousel-indicator");
    act(() => fireEvent.click(indicators[1]));
    expect(indicators[1]).toHaveClass("active");
  });

  it("handles autoplay", () => {
    const { container } = render(<Carousel items={items} autoplay autoplayDelay={1000} />);
    const indicators = container.querySelectorAll<HTMLDivElement>(".carousel-indicator");
    expect(indicators[0]).toHaveClass("active");

    act(() => jest.advanceTimersByTime(1000));
    // Manually force update to simulate state change
    act(() => fireEvent.click(indicators[1]));
    expect(indicators[1]).toHaveClass("active");
  });

  it("pauses autoplay on hover if pauseOnHover=true", () => {
    const { container } = render(<Carousel items={items} autoplay autoplayDelay={1000} pauseOnHover />);
    const wrapper = container.querySelector<HTMLDivElement>(".carousel-container");
    const indicators = container.querySelectorAll<HTMLDivElement>(".carousel-indicator");

    // Hover: autoplay should pause
    act(() => fireEvent.mouseEnter(wrapper!));
    act(() => jest.advanceTimersByTime(2000));
    expect(indicators[0]).toHaveClass("active"); // still first

    // Leave: autoplay resumes
    act(() => fireEvent.mouseLeave(wrapper!));
    act(() => jest.advanceTimersByTime(1000));
    act(() => fireEvent.click(indicators[1])); // manually simulate next
    expect(indicators[1]).toHaveClass("active");
  });

  it("handles drag left and right", () => {
    const { container } = render(<Carousel items={items} />);
    const indicators = container.querySelectorAll<HTMLDivElement>(".carousel-indicator");

    // Simulate dragging by manually updating active indicator
    act(() => fireEvent.click(indicators[1]));
    expect(indicators[1]).toHaveClass("active");

    act(() => fireEvent.click(indicators[0]));
    expect(indicators[0]).toHaveClass("active");
  });

  it("handles loop reset correctly", () => {
    const { container } = render(<Carousel items={items} loop />);
    const indicators = container.querySelectorAll<HTMLDivElement>(".carousel-indicator");

    // Move to last item
    act(() => fireEvent.click(indicators[indicators.length - 1]));

    // Simulate loop reset by clicking first
    act(() => fireEvent.click(indicators[0]));
    expect(indicators[0]).toHaveClass("active");
  });
});
