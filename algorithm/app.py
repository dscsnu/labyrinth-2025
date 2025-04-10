import pygame
import os
import time
import ast
import shutil

os.chdir(os.path.dirname(__file__))

pygame.init()

# --- Initial Setup ---
WIDTH, HEIGHT = 1280, 720
WIN = pygame.display.set_mode((WIDTH, HEIGHT), pygame.RESIZABLE)
pygame.display.set_caption("Image Approval System")

# --- Paths ---
LEFT_IMAGE_FOLDER = os.path.abspath("patterns_output")
RIGHT_IMAGE_FOLDER = os.path.abspath("spells")
OUTPUT_FOLDER = os.path.abspath("sorted_output")
REJECTED_FOLDER = os.path.abspath("rejected")

# --- Colors ---
WHITE = (255, 255, 255)
BLACK = (0, 0, 0)
GRAY = (200, 200, 200)
BLUE = (70, 130, 180)
DARK_BLUE = (50, 90, 130)
RED = (200, 60, 60)
DARK_RED = (150, 40, 40)
GREEN = (0, 255, 0)
BG_COLOR = (30, 41, 59)
DIVIDER_COLOR = (70, 80, 100)
ARROW_BG = (50, 50, 50)
ARROW_SHADOW = (0, 0, 0, 60)

# --- Fonts ---
FONT = pygame.font.SysFont("arial", 24)
BIG_FONT = pygame.font.SysFont("arial", 32, bold=True)

# --- Create folders if they don't exist ---
os.makedirs(OUTPUT_FOLDER, exist_ok=True)
os.makedirs(REJECTED_FOLDER, exist_ok=True)

# --- Load Images ---
def load_images(folder):
    """Load image file paths without loading the actual images yet"""
    files = sorted([f for f in os.listdir(folder) if f.lower().endswith(('.png', '.jpg', '.jpeg'))])
    paths = [os.path.join(folder, f) for f in files]
    return paths, files

def load_image(path, size=None):
    """Load a single image with optional resizing"""
    img = pygame.image.load(path).convert_alpha()
    if size:
        img = pygame.transform.smoothscale(img, size)
    return img

# Only store paths, not loaded images
left_paths, left_files = load_images(LEFT_IMAGE_FOLDER)
right_paths, right_files = load_images(RIGHT_IMAGE_FOLDER)
left_size = len(left_paths)

with open("patterns_nodes.txt", "r") as f:
    content = f.read()

left_nodelist = ast.literal_eval(content)

# Cache for loaded images
image_cache = {}
image_cache_size = 10  # Number of images to keep in cache

left_index = 0
right_index = 0
previous_index = 0
action_stack = []  # Tracks actions (approval/rejection) for back functionality
node_stack = []  # Tracks nodes for back functionality

# --- Rectangles ---
approve_button = pygame.Rect(0, 0, 160, 50)
reject_button = pygame.Rect(0, 0, 160, 50)
back_button = pygame.Rect(0, 0, 160, 50)
left_arrow = pygame.Rect(0, 0, 60, 60)
right_arrow = pygame.Rect(0, 0, 60, 60)

# --- State ---
approval_time = 0
show_message = False
message_text = ""
message_color = GREEN
last_switch_time = 0
image_switch_cooldown = 0.25
message_duration = 1.5  # How long to show messages in seconds

# --- Hold Scroll ---
holding_left = False
holding_right = False
hold_interval = 0.15
next_hold_time = 0

# --- Keyboard Hold ---
key_held_left = False
key_held_right = False

# --- UI Draw ---
def draw_button(rect, text, color, hover=False):
    offset = 4
    draw_rect = pygame.Rect(rect.x, rect.y + (2 if hover else 0), rect.width, rect.height)
    if hover:
        shadow_rect = pygame.Rect(draw_rect.x + offset, draw_rect.y + offset, draw_rect.width, draw_rect.height)
        pygame.draw.rect(WIN, BLACK, shadow_rect, border_radius=12)
    pygame.draw.rect(WIN, color, draw_rect, border_radius=12)
    label = FONT.render(text, True, WHITE)
    WIN.blit(label, (draw_rect.x + (rect.width - label.get_width()) // 2, draw_rect.y + 12))

def draw_arrow(center_x, center_y, direction="left", hover=False):
    if hover:
        shadow_surf = pygame.Surface((72, 72), pygame.SRCALPHA)
        pygame.draw.circle(shadow_surf, ARROW_SHADOW, (36, 36), 30)
        WIN.blit(shadow_surf, (center_x - 36 + 2, center_y - 36 + 2))
    pygame.draw.circle(WIN, ARROW_BG, (center_x, center_y), 30)
    if direction == "left":
        points = [(center_x + 10, center_y - 15), (center_x - 10, center_y), (center_x + 10, center_y + 15)]
    else:
        points = [(center_x - 10, center_y - 15), (center_x + 10, center_y), (center_x - 10, center_y + 15)]
    pygame.draw.polygon(WIN, WHITE, points)

def get_cached_image(path, size):
    """Get an image from cache or load it if not cached"""
    cache_key = (path, size[0], size[1])
    if cache_key not in image_cache:
        # If cache is full, remove the oldest entry
        if len(image_cache) >= image_cache_size:
            oldest_key = next(iter(image_cache))
            del image_cache[oldest_key]
        # Load and cache the image
        image_cache[cache_key] = load_image(path, size)
    return image_cache[cache_key]

def draw():
    global show_message
    WIN.fill(BG_COLOR)
    screen_width, screen_height = WIN.get_size()
    pygame.draw.line(WIN, DIVIDER_COLOR, (screen_width // 2, 0), (screen_width // 2, screen_height), 3)

    image_max_size = min(screen_width // 3, screen_height // 2)
    image_size = min(image_max_size, 350)
    img_y = screen_height // 2 - image_size // 2 - 40
    left_img_x = screen_width // 4 - image_size // 2
    right_img_x = 3 * screen_width // 4 - image_size // 2

    # Display "No images" message if left folder is empty
    if len(left_paths) == 0:
        no_images_text = FONT.render("No images in source folder", True, WHITE)
        WIN.blit(no_images_text, (left_img_x + (image_size - no_images_text.get_width()) // 2, 
                                 img_y + image_size // 2))
    # Draw left image if available
    elif left_index < len(left_paths):
        left_img = get_cached_image(left_paths[left_index], (image_size, image_size))
        left_rect = WIN.blit(left_img, (left_img_x, img_y))
        left_count = FONT.render(f"{left_size - len(left_paths) + 1}/{left_size}", True, WHITE)
        count_rect = left_count.get_rect(center=(left_rect.centerx, img_y + image_size + 20))
        WIN.blit(left_count, count_rect)

    # Display message if right folder is empty
    if len(right_paths) == 0:
        no_images_text = FONT.render("No images in right folder", True, WHITE)
        WIN.blit(no_images_text, (right_img_x + (image_size - no_images_text.get_width()) // 2, 
                                 img_y + image_size // 2))
    # Draw right image if available
    elif right_paths:
        right_path_index = right_index % len(right_paths)
        right_img = get_cached_image(right_paths[right_path_index], (image_size, image_size))
        right_rect = WIN.blit(right_img, (right_img_x, img_y))
        right_count = FONT.render(f"{right_path_index + 1}/{len(right_paths)}", True, WHITE)
        count_rect = right_count.get_rect(center=(right_rect.centerx, img_y + image_size + 20))
        WIN.blit(right_count, count_rect)

    arrow_y = img_y + image_size // 2
    mouse_pos = pygame.mouse.get_pos()
    left_arrow.center = (right_img_x - 70, arrow_y)
    right_arrow.center = (right_img_x + image_size + 70, arrow_y)

    draw_arrow(*left_arrow.center, "left", hover=left_arrow.collidepoint(mouse_pos))
    draw_arrow(*right_arrow.center, "right", hover=right_arrow.collidepoint(mouse_pos))

    button_y = img_y + image_size + 70
    buttons_center = right_img_x + image_size // 2
    total_width = approve_button.width + reject_button.width + 20

    approve_button.topleft = (buttons_center - total_width // 2, button_y)
    reject_button.topleft = (approve_button.right + 20, button_y)
    back_button.topleft = (buttons_center - back_button.width // 2, button_y + 60)

    # Display buttons with appropriate hover effects
    draw_button(approve_button, "APPROVE", DARK_BLUE if approve_button.collidepoint(mouse_pos) else BLUE, hover=approve_button.collidepoint(mouse_pos))
    draw_button(reject_button, "REJECT", DARK_RED if reject_button.collidepoint(mouse_pos) else RED, hover=reject_button.collidepoint(mouse_pos))
    draw_button(back_button, "BACK", GRAY, hover=back_button.collidepoint(mouse_pos))

    # Show message if needed
    if show_message:
        msg = BIG_FONT.render(message_text, True, message_color)
        WIN.blit(msg, (screen_width - msg.get_width() - 40, 30))

    pygame.display.update()

def reload_images():
    """Reload image paths after file operations"""
    global left_paths, left_files, image_cache
    left_paths, left_files = load_images(LEFT_IMAGE_FOLDER)
    # Clear cache when reloading images
    image_cache.clear()

def approve_current(now):
    global left_index, previous_index, approval_time, show_message, message_text, message_color

    if left_index >= len(left_paths) or len(left_paths) == 0:
        return

    left_filename = left_files[left_index]
    left_path = left_paths[left_index]
    left_nodes = left_nodelist[left_index]
    left_nodelist.pop(0)
    
    if len(right_files) == 0:
        return
        
    right_name = os.path.splitext(right_files[right_index % len(right_files)])[0]
    target_dir = os.path.join(OUTPUT_FOLDER, right_name)
    os.makedirs(target_dir, exist_ok=True)
    target_path = os.path.join(target_dir, left_filename)
    
    # Move instead of copy
    shutil.move(left_path, target_path)
    action_stack.append(("approved", left_index, target_path, left_path))
    
    # prints where the nodes are being assigned
    node_stack.append(left_nodes)
    print(f"{left_nodes} assigned to {right_name}")
    
    # Remember the current index before reloading
    current_index = left_index
    
    # After moving, reload the image list
    reload_images()
    
    # Stay at the same index after removing, unless we're at the end
    if current_index >= len(left_paths):
        left_index = max(0, len(left_paths) - 1)
    else:
        left_index = current_index
        
    approval_time = now
    message_text = "Approved!"
    message_color = GREEN
    show_message = True

def reject_current(now):
    global left_index, previous_index, approval_time, show_message, message_text, message_color
    
    if left_index >= len(left_paths) or len(left_paths) == 0:
        return
        
    left_filename = left_files[left_index]
    left_path = left_paths[left_index]
    target_path = os.path.join(REJECTED_FOLDER, left_filename)
    left_nodes = left_nodelist[left_index]
    left_nodelist.pop(0)
    
    # Move instead of copy
    shutil.move(left_path, target_path)
    action_stack.append(("rejected", left_index, target_path, left_path))

    # prints that the nodes were rejected
    node_stack.append(left_nodes)
    print(f"{left_nodes} rejected")
    
    # Remember the current index before reloading
    current_index = left_index
    
    # After moving, reload the image list
    reload_images()
    
    # Stay at the same index after removing, unless we're at the end
    if current_index >= len(left_paths):
        left_index = max(0, len(left_paths) - 1)
    else:
        left_index = current_index
        
    approval_time = now
    message_text = "Rejected!"
    message_color = RED
    show_message = True

def go_back(now):
    global left_index, approval_time, show_message, message_text, message_color
    if action_stack:
        action, index, target_path, original_path = action_stack.pop()

        #adds the nodes back to the original list
        left_nodelist.insert(0, node_stack.pop())

        print(f"{left_nodelist[0]} was retrieved")
        if os.path.exists(target_path):
            # Make sure the directory exists
            os.makedirs(os.path.dirname(original_path), exist_ok=True)
            # Move the file back to its original location
            shutil.move(target_path, original_path)
            
            # Reload images after moving it back
            reload_images()
            
            # Find the index of the restored file
            try:
                restored_filename = os.path.basename(original_path)
                restored_index = left_files.index(restored_filename)
                left_index = restored_index
            except ValueError:
                # If for some reason we can't find the file, just stay at current index
                pass
                
            approval_time = now
            message_text = "Undone"
            message_color = GRAY
            show_message = True

# --- Main Loop ---
def main():
    global right_index, holding_left, holding_right, next_hold_time
    global key_held_left, key_held_right, show_message

    clock = pygame.time.Clock()
    running = True
    fullscreen = False
    show_message = False

    while running:
        now = time.time()
        
        # Hide message after duration
        if show_message and now - approval_time > message_duration:
            show_message = False
        
        # Handle key holding for scrolling
        if right_paths:
            if (holding_left or key_held_left) and now >= next_hold_time:
                right_index = (right_index - 1) % len(right_paths)
                next_hold_time = now + hold_interval
            elif (holding_right or key_held_right) and now >= next_hold_time:
                right_index = (right_index + 1) % len(right_paths)
                next_hold_time = now + hold_interval

        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False

            elif event.type == pygame.KEYDOWN:
                if event.key == pygame.K_ESCAPE:
                    running = False
                elif event.key == pygame.K_RETURN:
                    approve_current(now)
                elif event.key == pygame.K_LEFT and right_paths:
                    key_held_left = True
                    next_hold_time = now
                elif event.key == pygame.K_RIGHT and right_paths:
                    key_held_right = True
                    next_hold_time = now
                elif event.key == pygame.K_f:
                    fullscreen = not fullscreen
                    pygame.display.set_mode((WIDTH, HEIGHT), pygame.FULLSCREEN if fullscreen else pygame.RESIZABLE)

            elif event.type == pygame.KEYUP:
                if event.key == pygame.K_LEFT:
                    key_held_left = False
                elif event.key == pygame.K_RIGHT:
                    key_held_right = False

            elif event.type == pygame.MOUSEBUTTONDOWN:
                if left_arrow.collidepoint(event.pos) and right_paths:
                    holding_left = True
                    next_hold_time = now
                elif right_arrow.collidepoint(event.pos) and right_paths:
                    holding_right = True
                    next_hold_time = now
                elif approve_button.collidepoint(event.pos):
                    approve_current(now)
                elif reject_button.collidepoint(event.pos):
                    reject_current(now)
                elif back_button.collidepoint(event.pos):
                    go_back(now)

            elif event.type == pygame.MOUSEBUTTONUP:
                holding_left = holding_right = False

        # Draw everything
        draw()
        
        # Cap the frame rate
        clock.tick(60)

    pygame.quit()

if __name__ == "__main__":
    main()