import sys
import time
from typing import Optional


class ProgressBar:

    def __init__(self, total: int, description: str = "", width: int = 50):

        self.total = total
        self.description = description
        self.width = width
        self.current = 0
        self.start_time = time.time()

    def update(self, step: int = 1):
        self.current = min(self.current + step, self.total)
        self._render()

    def set_progress(self, current: int):
        self.current = min(current, self.total)
        self._render()

    def _render(self):
        if self.total == 0:
            percent = 100
        else:
            percent = (self.current / self.total) * 100

        filled = int(self.width * self.current // self.total)
        bar = '█' * filled + '░' * (self.width - filled)

        elapsed = time.time() - self.start_time
        if self.current > 0 and self.current < self.total:
            eta = (elapsed / self.current) * (self.total - self.current)
            eta_str = f"ETA: {eta:.1f}s"
        elif self.current == self.total:
            eta_str = f"Done in {elapsed:.1f}s"
        else:
            eta_str = ""

        line = f"\r{self.description} [{bar}] {percent:>5.1f}% {eta_str}"

        sys.stdout.write(line)
        sys.stdout.flush()

        if self.current >= self.total:
            sys.stdout.write("\n")
            sys.stdout.flush()

    def finish(self):
        """Mark progress as complete"""
        self.current = self.total
        self._render()


class Spinner:

    FRAMES = ['⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏']

    def __init__(self, description: str = ""):

        self.description = description
        self.frame_index = 0
        self.start_time = time.time()
        self.running = False

    def update(self):
        if not self.running:
            self.running = True

        elapsed = time.time() - self.start_time
        frame = self.FRAMES[self.frame_index % len(self.FRAMES)]
        self.frame_index += 1

        line = f"\r{frame} {self.description} ({elapsed:.1f}s)"
        sys.stdout.write(line)
        sys.stdout.flush()

    def finish(self, message: Optional[str] = None):

        elapsed = time.time() - self.start_time

        if message:
            line = f"\r✓ {message} ({elapsed:.1f}s)\n"
        else:
            line = f"\r✓ {self.description} ({elapsed:.1f}s)\n"

        sys.stdout.write(line)
        sys.stdout.flush()
        self.running = False


def simulate_spinner(duration: float = 3.0, description: str = "Loading"):

    spinner = Spinner(description)
    start = time.time()

    while time.time() - start < duration:
        spinner.update()
        time.sleep(0.1)

    spinner.finish()


def simulate_progress_bar(steps: int = 100, description: str = "Processing"):

    bar = ProgressBar(steps, description)

    for i in range(steps):
        time.sleep(0.02)
        bar.update()


if __name__ == "__main__":
    """
    Testing Progress and Spinner indicators
    """
    print("Progress Indicators Demo\n")

    print("1. Spinner (indeterminate progress):")
    simulate_spinner(2.0, "Loading model")
    print()

    print("2. Progress Bar (determinate progress):")
    simulate_progress_bar(100, "Processing data")
    print()

    print("3. Multiple spinners:")
    simulate_spinner(1.5, "Downloading weights")
    simulate_spinner(1.0, "Initializing model")
    simulate_spinner(0.8, "Warming up")