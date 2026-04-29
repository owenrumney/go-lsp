function slugify(text) {
  return text
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9\s-]/g, "")
    .replace(/\s+/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-|-$/g, "");
}

document.addEventListener("DOMContentLoaded", () => {
  const seen = new Map();

  document.querySelectorAll("main h1, main h2, main h3").forEach((heading) => {
    const base = slugify(heading.textContent || "section") || "section";
    const count = seen.get(base) || 0;
    seen.set(base, count + 1);

    const id = count === 0 ? base : `${base}-${count + 1}`;
    heading.id = heading.id || id;

    const anchor = document.createElement("a");
    anchor.className = "heading-anchor";
    anchor.href = `#${heading.id}`;
    anchor.setAttribute("aria-label", `Permalink to ${heading.textContent}`);
    anchor.textContent = "#";
    heading.appendChild(anchor);
  });
});
