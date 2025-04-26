function formatDoctorOutput(raw) {
  if (typeof raw !== "string") return "";

  // Remove starting " if exists
  if (raw.startsWith('"')) {
    raw = raw.slice(1);
  }

  // Remove trailing " if exists (before trimming whitespace)
  if (raw.trimEnd().endsWith('"')) {
    raw = raw.trimEnd();
    raw = raw.slice(0, -1);
  }

  // Unescape escaped characters (like \n, \")
  raw = raw.replace(/\\"/g, '"').replace(/\\n/g, "\n");

  // Escape HTML to prevent XSS
  const escaped = raw
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");

  // Color special symbols and format line breaks
  return escaped
    .replace(/\[✓]/g, '<span class="symbol-green">[✓]</span>')
    .replace(/\[✗]/g, '<span class="symbol-red">[✗]</span>')
    .replace(/\[!]/g, '<span class="symbol-orange">[!]</span>')
    .replace(/\n/g, "<br>");
}



document.addEventListener("DOMContentLoaded", () => {
  loadFlutterDoctor(); // Loads local first inside this function
  loadPackageDiff();
  document.getElementById("sync-btn").addEventListener("click", syncPackages);
});

// Tab switch logic
document.querySelectorAll(".tab-btn").forEach((btn) => {
  btn.addEventListener("click", () => {
    document
      .querySelectorAll(".tab-btn")
      .forEach((b) => b.classList.remove("active"));
    btn.classList.add("active");

    document.querySelectorAll(".tab-content").forEach((tab) => {
      tab.style.display = "none";
    });

    const selectedTab = document.getElementById(btn.getAttribute("data-tab"));
    if (selectedTab) selectedTab.style.display = "block";
  });
});

// Attach the event listener to the Refresh button
document
  .getElementById("refresh-btn")
  .addEventListener("click", refreshPackageTabs);

const themeBtn = document.getElementById("theme-toggle-btn");
const themeLink = document.getElementById("theme-style");

themeBtn.addEventListener("click", () => {
  const isDark = themeLink.getAttribute("href").includes("style_dark.css");

  if (isDark) {
    themeLink.setAttribute("href", "/views/assets/css/style_light.css");
    themeBtn.textContent = "☀️ Light Mode";
  } else {
    themeLink.setAttribute("href", "/views/assets/css/style_dark.css");
    themeBtn.textContent = "🌙 Dark Mode";
  }
});
