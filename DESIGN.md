---
name: AutoSSL Console
description: A quiet operational console for issuing and distributing domain certificates.
colors:
  background: "hsl(210 20% 98%)"
  foreground: "hsl(150 18% 13%)"
  card: "hsl(0 0% 100%)"
  primary: "hsl(158 62% 27%)"
  primary-foreground: "hsl(0 0% 100%)"
  secondary: "hsl(150 12% 94%)"
  muted-foreground: "hsl(150 6% 42%)"
  border: "hsl(150 12% 86%)"
  destructive: "hsl(0 68% 47%)"
  ring: "hsl(158 62% 32%)"
typography:
  headline:
    fontFamily: "Inter, ui-sans-serif, system-ui, sans-serif"
    fontSize: "24px"
    fontWeight: 600
    lineHeight: 1.25
    letterSpacing: "0"
  title:
    fontFamily: "Inter, ui-sans-serif, system-ui, sans-serif"
    fontSize: "18px"
    fontWeight: 600
    lineHeight: 1.4
    letterSpacing: "0"
  body:
    fontFamily: "Inter, ui-sans-serif, system-ui, sans-serif"
    fontSize: "14px"
    fontWeight: 400
    lineHeight: 1.5
    letterSpacing: "0"
  label:
    fontFamily: "Inter, ui-sans-serif, system-ui, sans-serif"
    fontSize: "12px"
    fontWeight: 500
    lineHeight: 1.4
    letterSpacing: "0"
rounded:
  sm: "6px"
  md: "8px"
  lg: "10px"
  full: "9999px"
spacing:
  xs: "4px"
  sm: "8px"
  md: "16px"
  lg: "24px"
  xl: "32px"
components:
  button-primary:
    backgroundColor: "{colors.primary}"
    textColor: "{colors.primary-foreground}"
    rounded: "{rounded.md}"
    padding: "0 16px"
    height: "44px"
  button-outline:
    backgroundColor: "{colors.card}"
    textColor: "{colors.foreground}"
    rounded: "{rounded.md}"
    padding: "0 16px"
    height: "44px"
  input:
    backgroundColor: "{colors.card}"
    textColor: "{colors.foreground}"
    rounded: "{rounded.md}"
    padding: "8px 12px"
    height: "44px"
  card:
    backgroundColor: "{colors.card}"
    textColor: "{colors.foreground}"
    rounded: "{rounded.lg}"
    padding: "16px"
---

# Design System: AutoSSL Console

## Overview

**Creative North Star: "The Quiet Operations Desk"**

AutoSSL is a compact administrative surface where clarity outranks decoration. It uses a cool near-white workspace, crisp white work surfaces, and a restrained green accent to make certificate actions easy to scan and verify.

The interface should feel calm, trustworthy, and direct. Hierarchy comes from spacing, typography, borders, and a small amount of ambient shadow rather than oversized headings or decorative panels.

**Key Characteristics:**

- Compact operational density with generous touch targets.
- Green is reserved for identity, focus, and primary actions.
- Certificate URLs remain visible, selectable, and paired with explicit copy actions.
- Destructive actions are red, secondary, and always confirmed.

## Colors

The palette combines cool neutrals with an operational green and a narrowly scoped destructive red.

### Primary

- **Operational Green** (`primary`): brand mark, primary actions, selection, and focus emphasis.

### Neutral

- **Cool Workspace** (`background`): page canvas.
- **Clear Surface** (`card`): forms, tables, cards, and dialogs.
- **Ink Green** (`foreground`): primary text and icons.
- **Quiet Fill** (`secondary`): subdued controls, table headers, and hover states.
- **Reference Gray** (`muted-foreground`): labels and supporting text.
- **Soft Divider** (`border`): structural separation without visual weight.
- **Action Red** (`destructive`): deletion and destructive confirmation only.

**The One Accent Rule.** Operational green is the only positive accent; do not introduce competing brand colors.

## Typography

**Display Font:** Inter with the system sans-serif stack
**Body Font:** Inter with the system sans-serif stack
**Label/Mono Font:** The system monospace stack is used only for certificate URLs

**Character:** Neutral, compact, and readable. Weight and whitespace carry hierarchy; letter spacing remains normal.

### Hierarchy

- **Headline** (semibold, 24px, 1.25): page titles only.
- **Title** (semibold, 18px, 1.4): product identity, card domains, and dialog titles.
- **Body** (regular, 14px, 1.5): actions and operational content.
- **Label** (medium, 12px, 1.4): table headings, field labels, and supporting metadata.

**The Compact Type Rule.** Keep operational labels and table content small enough to scan while preserving readable line height and 44px controls.

## Layout

Content is centered in a 1280px maximum-width workspace with 16px mobile gutters, 24px small-screen gutters, and 32px large-screen gutters. A single 64px top bar holds identity and sign-out; the page title and primary action follow directly below.

Certificate rows use a four-column table at 1024px and above. Below 1024px, each certificate becomes a stacked card so both URLs, copy actions, and delete remain visible without horizontal scrolling. Vertical rhythm uses an 8px base and its multiples.

## Elevation & Depth

The system is flat by default. Borders separate table regions and card content; soft ambient shadows lift login, certificate, and dialog surfaces only when a bounded surface needs definition.

**The Flat Workspace Rule.** Do not wrap page sections in decorative cards. Reserve elevation for actual records, forms, and modal content.

## Shapes

Controls use gently curved 8px corners, while cards and dialogs use 10px corners. Count badges use a full pill. Borders are thin and quiet; shapes remain practical rather than playful.

## Components

### Buttons

- **Shape:** 8px corners and a stable 44px height.
- **Primary:** operational green with white text; used for login, adding domains, and affirmative actions.
- **Hover / Focus:** subtle color shift with a two-pixel green focus ring and offset.
- **Outline / Ghost:** white or transparent surfaces with quiet neutral hover fills.
- **Destructive:** red is reserved for delete actions and their confirmation.

### Cards / Containers

- **Corner Style:** gently curved (10px).
- **Background:** clear white surface against the cool workspace.
- **Shadow Strategy:** soft ambient shadow; no stacked or nested cards.
- **Border:** dividers appear inside certificate cards; outer borders are optional when shadow already defines the edge.
- **Internal Padding:** 16px for record content and 24-32px for login/dialog surfaces.

### Inputs / Fields

- **Style:** 44px high white fields with an 8px radius and quiet neutral border.
- **Focus:** darker green border plus a translucent green focus ring.
- **Error / Disabled:** errors use destructive red; disabled controls reduce opacity without losing labels.

### Navigation

The 64px white top bar contains the shield mark and AutoSSL wordmark on the left and a compact sign-out action on the right. It has no sidebar, breadcrumbs, or secondary navigation.

### Certificate Record

Every record leads with the domain, then exposes the certificate and private-key URLs in monospace. Each URL has its own copy action, while delete stays visually separate and requires confirmation.

## Do's and Don'ts

### Do:

- **Do** keep primary workflows on one page and make their current state obvious.
- **Do** preserve 44px control targets, keyboard focus, semantic labels, and reduced-motion behavior.
- **Do** switch records to cards below 1024px so all actions remain discoverable.
- **Do** show copyable certificate URLs; scripts are the intended consumer.

### Don't:

- **Don't** add certificate download buttons to the console.
- **Don't** expose API bearer tokens or call administrative Go endpoints from browser code.
- **Don't** add a sidebar, marketing hero, decorative gradients, or nested cards.
- **Don't** use green for deletion or red for routine status.
