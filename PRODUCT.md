# AutoSSL

<!-- impeccable:product-schema 1 -->

## Platform

web

## Stack

Go and Echo API with a React, React Router, Vite, shadcn/ui, Tailwind CSS, and Lucide web console. The web server acts as a BFF and is packaged with the Go service in one container.

## Users

The primary user is the administrator who issues and removes domain certificates and copies distribution links for deployment scripts.

## Product Purpose

AutoSSL centralizes SSL certificate issuance and exposes stable certificate and private-key URLs for automated distribution.

## Positioning

Certificate issuance, storage, renewal, and script-facing distribution remain in the existing Go service while the console provides a human-operated server-side interface.

## Operating Context

Administrators sign in, inspect the certificate list, issue a certificate for a base domain, copy certificate or private-key URLs, and delete certificates. Deployment uses one container.

## Capabilities and Constraints

- Browser requests for administrative operations terminate at the web server; the web server calls the Go API.
- The console has a login page and one certificate-management page.
- Certificate URLs are displayed for copying. The console does not provide download buttons.
- The Go service and Node web service ship in the same container.

## Evidence on Hand

The repository contains the working Go API, ACME integration, SQLite repository, JWT login, and certificate distribution routes. No separate brand assets or commercial claims are supplied.

## Product Principles

- Keep certificate operations explicit and easy to verify.
- Keep backend credentials and tokens out of browser JavaScript.
- Preserve stable script-facing certificate URLs.
- Prefer a compact operational interface over dashboard decoration.

## Accessibility & Inclusion

The console must support keyboard operation, visible focus, labeled controls, readable contrast, and reduced motion preferences.
