# 🥋 Go Katas 🥋

>  "I fear not the man who has practiced 10,000 kicks once, but I fear the man who has practiced one kick 10,000 times."
(Bruce Lee)

## What should it be?
- Go is simple to learn, but nuanced to master. The difference between "working code" and "idiomatic code" often lies in details such as safety, memory efficiency, and concurrency control.

- This repository is a collection of **Daily Katas**: small, standalone coding challenges designed to drill specific Go patterns into your muscle memory.

## What should it NOT be? 

- This is not intended to teach coding, having Go as the programming mean. Not even intended to teach you Go **in general**
- The focus should be as much as possible challenging oneself to solve common software engineering problems **the Go way**. 
- Several seasoned developers spent years learning and applying best-practices at prod-grade context. Once they decide to switch to go, they would face two challanges:
  - Is there a window of knowledge transform here, so that I don't have to through years of my career from the window at start from zero?
  - If yes, the which parts should I focus on to recognize the mismatches and use them the expected way in the Go land?

## How to Use This Repo
1.  **Pick a Kata:** Navigate to any `XX-kata-yy` folder.
2.  **Read the Challenge:** Open the `README.md` inside that folder. It defines the Goal, the Constraints, and the "Idiomatic Patterns" you must use.
3.  **Solve It:** Initialize a module inside the folder and write your solution.
4.  **Reflect:** Compare your solution with the provided "Reference Implementation" (if available) or the core patterns listed.

## Contribution Guidelines

### Have a favorite Go pattern?
1. Create a new folder `XX-your-topic`. (`XX` is an ordinal number)
2. Copy the [README_TEMPLATE.md](./README_TEMPLATE.md) to the new folder as `README.md`
3. Define the challenge: focus on **real-world scenarios** (e.g., handling timeouts, zero-allocation sets), and **idiomatic Go**, not just algorithmic puzzles.
4. **Optionally**, create a `main.go` or any other relevant files under the project containing blueprint of the implementation, **as long as you think it reduces confusion and keeps the implementation  focused**
5. Submit a PR.

### Using the script

You can use the shorthand script to add a new challenge, it will create a new folder and a new README.md file under it:
```bash
./add.sh my-very-creative-challange
```
This will create a new folder `21-my-very-creative-challange` (in case the latest challange was under the folder name `20-latest-name-here`) and add a `README.md` under it

```bash
medunes@medunes:~/projects/go-kata$  ls 21-my-very-creative-challange/
README.md
```
