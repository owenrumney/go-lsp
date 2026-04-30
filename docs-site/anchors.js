function slugify(text) {
  return text
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9\s-]/g, "")
    .replace(/\s+/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-|-$/g, "");
}

function appendGoTokens(code, source) {
  const keywords = new Set([
    "break", "case", "chan", "const", "continue", "default", "defer", "else",
    "fallthrough", "for", "func", "go", "goto", "if", "import", "interface",
    "map", "package", "range", "return", "select", "struct", "switch", "type",
    "var",
  ]);
  const builtins = new Set([
    "append", "bool", "byte", "cap", "close", "complex", "comparable", "copy",
    "delete", "error", "false", "float32", "float64", "imag", "int", "int16",
    "int32", "int64", "int8", "iota", "len", "make", "new", "nil", "panic",
    "print", "println", "real", "recover", "rune", "string", "true", "uint",
    "uint16", "uint32", "uint64", "uint8", "uintptr", "any",
  ]);
  const tokenPattern = /\/\/.*|"(?:\\.|[^"\\])*"|`[\s\S]*?`|'(?:\\.|[^'\\])*'|\b\d+(?:\.\d+)?\b|\b[A-Za-z_]\w*\b/g;

  let lastIndex = 0;

  function appendText(text) {
    if (text) {
      code.appendChild(document.createTextNode(text));
    }
  }

  function appendToken(text, className) {
    const span = document.createElement("span");
    span.className = className;
    span.textContent = text;
    code.appendChild(span);
  }

  source.replace(tokenPattern, (match, offset) => {
    appendText(source.slice(lastIndex, offset));

    if (match.startsWith("//")) {
      appendToken(match, "token comment");
    } else if (match.startsWith('"') || match.startsWith("'") || match.startsWith("`")) {
      appendToken(match, "token string");
    } else if (/^\d/.test(match)) {
      appendToken(match, "token number");
    } else if (keywords.has(match)) {
      appendToken(match, "token keyword");
    } else if (builtins.has(match)) {
      appendToken(match, "token builtin");
    } else {
      appendText(match);
    }

    lastIndex = offset + match.length;
    return match;
  });

  appendText(source.slice(lastIndex));
}

async function copyText(text) {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(text);
    return;
  }

  const input = document.createElement("textarea");
  input.value = text;
  input.setAttribute("readonly", "");
  input.style.position = "absolute";
  input.style.left = "-9999px";
  document.body.appendChild(input);
  input.select();
  document.execCommand("copy");
  input.remove();
}

function enhanceCodeBlocks() {
  document.querySelectorAll("pre > code").forEach((code) => {
    const source = code.textContent || "";
    const pre = code.parentElement;

    if (!pre) {
      return;
    }

    pre.classList.add("code-block");

    if (code.classList.contains("language-go")) {
      pre.dataset.language = "go";
      code.textContent = "";
      appendGoTokens(code, source);
    }

    const button = document.createElement("button");
    button.className = "code-copy";
    button.type = "button";
    button.textContent = "Copy";
    button.setAttribute("aria-label", "Copy code to clipboard");
    button.addEventListener("click", async () => {
      try {
        await copyText(source);
        button.textContent = "Copied";
        button.classList.add("copied");
        window.setTimeout(() => {
          button.textContent = "Copy";
          button.classList.remove("copied");
        }, 1800);
      } catch {
        button.textContent = "Failed";
        window.setTimeout(() => {
          button.textContent = "Copy";
        }, 1800);
      }
    });

    pre.appendChild(button);
  });
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

  const sidebar = document.querySelector(".sidebar");
  const toggle = document.querySelector(".menu-toggle");
  if (sidebar && toggle) {
    toggle.addEventListener("click", () => {
      const open = sidebar.classList.toggle("nav-open");
      toggle.setAttribute("aria-expanded", open ? "true" : "false");
    });
  }

  enhanceCodeBlocks();
});
