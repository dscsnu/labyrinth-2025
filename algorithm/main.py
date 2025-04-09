import matplotlib.pyplot as plt
import itertools
import os

# Make output folder
os.makedirs("patterns_output", exist_ok=True)

# Node No : (x,y)
positions = {
    0: (0, 2),
    1: (1, 2),
    2: (2, 2),
    3: (0, 1),
    4: (1, 5),
    5: (2, 1),
    6: (0, 0),
    7: (1, 0),
    8: (2, 0),
    9: (1, -1),
}

# Get all 3-vertex and 4-vertex combinations
combos_3 = list(itertools.combinations(range(10), 3))
combos_4 = list(itertools.combinations(range(10), 4))
# combos_8 = list(itertools.combinations(range(10), 8))

# Combine all combos
all_combos = combos_3 + combos_4 #+ combos_8

# Plot and save each one
for i, combo in enumerate(all_combos):
    print(combo)
    fig, ax = plt.subplots(figsize=(4, 4))
    ax.set_title(f"Nodes: {combo}")
    ax.set_xlim(-1, 3)
    ax.set_ylim(-2, 3)
    ax.set_aspect('equal')
    ax.axis('off')

    for node, (x, y) in positions.items():
        color = 'gray' if node in combo else 'lightblue'
        size = 300 if node in combo else 200
        ax.scatter(x, y, s=size, color=color, edgecolors='black', zorder=3)
        ax.text(x, y, str(node), ha='center', va='center', zorder=4, fontsize=9)

    plt.tight_layout()
    plt.savefig(f"patterns_output/pattern_{i:03d}.png")
    plt.close()
