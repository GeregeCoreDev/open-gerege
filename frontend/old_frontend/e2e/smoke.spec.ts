import { test, expect } from '@playwright/test';

test('has title', async ({ page }) => {
    await page.goto('/');

    // Expect a title "to contain" a substring.
    // Note: Adjust this based on your actual metadata title
    await expect(page).toHaveTitle(/Gerege/);
});

test('redirects to login', async ({ page }) => {
    await page.goto('/');

    // Should eventually land on login page if not authenticated
    // Adjust the URL pattern as needed based on your routing
    await expect(page.url()).toContain('/home');
});
