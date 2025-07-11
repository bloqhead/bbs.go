# Color Palette - "The Naked Gun" Movie Poster

This CSS file contains a color palette extracted from "The Naked Gun" movie poster, featuring vibrant blues, striking magentas, and sophisticated neutrals.

## Color Variables

### Primary Colors
- `--sky-blue`: #87CEEB - Light blue from the poster's sky background
- `--bright-magenta`: #E91E63 - Bold pink from the movie title
- `--deep-pink`: #C2185B - Darker variant of the magenta

### Neutral Colors
- `--charcoal-gray`: #2C2C2C - Dark gray from the suit
- `--medium-gray`: #6B6B6B - Mid-tone gray
- `--light-gray`: #A8A8A8 - Light gray for subtle elements
- `--off-white`: #F5F5F5 - Soft white background
- `--pure-white`: #FFFFFF - Pure white
- `--pure-black`: #000000 - Pure black

### Accent Colors
- `--flesh-tone`: #F4C2A1 - Warm skin tone
- `--burgundy-red`: #722F37 - Deep red from the tie
- `--warm-brown`: #8B4513 - Brown accent

### Semantic Colors
- `--primary`: Bright magenta (main brand color)
- `--secondary`: Sky blue (secondary brand color)
- `--accent`: Burgundy red (accent color)
- `--background`: Off-white (page background)
- `--surface`: Pure white (card/content backgrounds)
- `--text-primary`: Charcoal gray (main text)
- `--text-secondary`: Medium gray (secondary text)
- `--text-inverse`: Pure white (text on dark backgrounds)

## Utility Classes

### Background Colors
```css
.bg-primary, .bg-secondary, .bg-accent
.bg-sky-blue, .bg-magenta, .bg-charcoal
.bg-light-gray, .bg-burgundy
```

### Text Colors
```css
.text-primary, .text-secondary, .text-inverse
.text-magenta, .text-sky-blue, .text-burgundy
.text-white, .text-black
```

### Border Colors
```css
.border-primary, .border-secondary, .border-accent
.border-magenta, .border-sky-blue, .border-gray
```

### Button Components
- `.btn-primary` - Magenta button with white text
- `.btn-secondary` - Sky blue button with dark text
- `.btn-outline` - Transparent button with magenta border

### Card Components
- `.card` - Basic white card with shadow
- `.card-accent` - Card with magenta left border

### Gradient Utilities
- `.gradient-sky` - Sky blue to white gradient
- `.gradient-magenta` - Magenta to deep pink gradient
- `.gradient-hero` - Sky blue to magenta gradient

### Shadow Utilities
- `.shadow-sm`, `.shadow-md`, `.shadow-lg` - Various shadow sizes
- `.shadow-magenta` - Magenta-colored shadow

### Layout Utilities
- `.container` - Responsive container (max-width: 1200px)
- `.hero-section` - Hero section with gradient background

### Typography
- `.heading-primary` - Magenta headings
- `.heading-accent` - Burgundy headings
- `.text-muted` - Gray muted text

## Usage

1. Link the CSS file in your HTML:
```html
<link rel="stylesheet" href="color-palette.css">
```

2. Use CSS variables directly:
```css
.my-element {
    background-color: var(--primary);
    color: var(--text-inverse);
}
```

3. Use utility classes:
```html
<div class="bg-primary text-inverse">
    <h1 class="heading-primary">Hello World</h1>
    <button class="btn-outline">Click me</button>
</div>
```

## Example

See `example.html` for a complete demonstration of all colors and utility classes in action.

## Color Accessibility

This palette maintains good contrast ratios:
- Dark text on light backgrounds
- White text on dark/colored backgrounds
- Medium gray for secondary text elements

## Browser Support

CSS custom properties are supported in all modern browsers (IE 11+ with partial support).