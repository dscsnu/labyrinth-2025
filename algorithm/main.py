import pygame

pygame.init()
screen = pygame.display.set_mode((1400, 800))
pygame.display.set_caption("Image Approval System")

# Colors
WHITE = (255, 255, 255)
DARK = (15, 30, 50)
GRAY = (30, 30, 30)
BLUE = (70, 140, 200)
RED = (200, 70, 70)
SHADOW = (0, 0, 0)

# Load example image
image = pygame.image.load("image.png")

# # Dummy thumbnails list
thumbnails = [image] * 7  # Replace with actual loading logic
current_index = 0
pattern_index = 0
total_patterns = 330

# Fonts
font = pygame.font.SysFont("arial", 24)
big_font = pygame.font.SysFont("arial", 36)

clock = pygame.time.Clock()
hold_left = hold_right = False
hold_timer = 0

def draw_arrow_button(x, y, direction="left"):
    # Shadow
    pygame.draw.circle(screen, SHADOW, (x+4, y+4), 30)
    # Button
    pygame.draw.circle(screen, GRAY, (x, y), 30)
    # Arrow
    if direction == "left":
        pygame.draw.polygon(screen, WHITE, [(x+8, y-15), (x-10, y), (x+8, y+15)])
    else:
        pygame.draw.polygon(screen, WHITE, [(x-8, y-15), (x+10, y), (x-8, y+15)])

def draw_button(x, y, w, h, text, color):
    pygame.draw.rect(screen, SHADOW, (x+4, y+4, w, h), border_radius=10)
    pygame.draw.rect(screen, color, (x, y, w, h), border_radius=10)
    label = big_font.render(text, True, WHITE)
    screen.blit(label, (x + w//2 - label.get_width()//2, y + h//2 - label.get_height()//2))

running = True
while running:
    screen.fill(DARK)

    # Left panel (placeholder for graph)
    pygame.draw.rect(screen, WHITE, (150, 200, 300, 300))
    graph_text = font.render(f"{pattern_index+1}/{total_patterns}", True, WHITE)
    screen.blit(graph_text, (250, 520))

    # Right panel - current image
    # image = thumbnails[current_index]
    image_rect = image.get_rect(center=(1000, 300))
    pygame.draw.rect(screen, WHITE, image_rect.inflate(20, 20), border_radius=15)
    screen.blit(image, image_rect)

    # Image index
    idx_text = font.render(f"{current_index+1}/{len(thumbnails)}", True, WHITE)
    screen.blit(idx_text, (1000 - idx_text.get_width()//2, 450))

    # Arrows
    draw_arrow_button(880, 300, "left")
    draw_arrow_button(1120, 300, "right")

    # Approve/Reject buttons (centered under image)
    draw_button(930, 500, 120, 60, "APPROVE", BLUE)
    draw_button(1070, 500, 120, 60, "REJECT", RED)

    pygame.display.flip()
    clock.tick(60)

    # Key hold
    keys = pygame.key.get_pressed()
    if keys[pygame.K_LEFT]:
        hold_left = True
        hold_timer += 1
        if hold_timer % 7 == 0:
            current_index = (current_index - 1) % len(thumbnails)
    elif keys[pygame.K_RIGHT]:
        hold_right = True
        hold_timer += 1
        if hold_timer % 7 == 0:
            current_index = (current_index + 1) % len(thumbnails)
    else:
        hold_timer = 0

    # Events
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            running = False

        elif event.type == pygame.MOUSEBUTTONDOWN:
            mx, my = pygame.mouse.get_pos()

            # Left arrow
            if (mx - 880)**2 + (my - 300)**2 <= 30**2:
                current_index = (current_index - 1) % len(thumbnails)

            # Right arrow
            elif (mx - 1120)**2 + (my - 300)**2 <= 30**2:
                current_index = (current_index + 1) % len(thumbnails)

            # Approve
            elif 930 <= mx <= 1050 and 500 <= my <= 560:
                print(f"Approved image {current_index+1}")

            # Reject
            elif 1070 <= mx <= 1190 and 500 <= my <= 560:
                print(f"Rejected image {current_index+1}")

        elif event.type == pygame.KEYDOWN:
            if event.key == pygame.K_LEFT:
                current_index = (current_index - 1) % len(thumbnails)
            elif event.key == pygame.K_RIGHT:
                current_index = (current_index + 1) % len(thumbnails)

pygame.quit()
