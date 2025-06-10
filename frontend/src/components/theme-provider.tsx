import { createContext, useContext, useEffect, useState } from "react";

type Theme = "dark" | "light" | "system";

type ThemeProviderProps = {
  children: React.ReactNode;
  defaultTheme?: Theme;
  storageKey?: string;
};

type ThemeProviderState = {
  theme: Theme;
  setTheme: (theme: Theme) => void;
};

const initialState: ThemeProviderState = {
  theme: "system",
  setTheme: () => null,
};

const ThemeProviderContext = createContext<ThemeProviderState>(initialState);

export function ThemeProvider({
  children,
  defaultTheme = "system",
  storageKey = "vite-ui-theme",
  ...props
}: ThemeProviderProps) {
  const [theme, setTheme] = useState<Theme>(
    () => (localStorage.getItem(storageKey) as Theme) || defaultTheme
  );

  useEffect(() => {
    const root = window.document.documentElement;

    root.classList.remove("light", "dark");

    if (theme === "system") {
      const systemTheme = window.matchMedia("(prefers-color-scheme: dark)")
        .matches
        ? "dark"
        : "light";

      root.classList.add(systemTheme);
      return;
    }

    root.classList.add(theme);
  }, [theme]);

  const value = {
    theme,
    setTheme: (theme: Theme) => {
      localStorage.setItem(storageKey, theme);
      setTheme(theme);
    },
  };

  return (
    <ThemeProviderContext.Provider {...props} value={value}>
      {/*
        This is the main application container.
        - min-h-screen: Ensures the container takes at least the full viewport height.
        - bg-black: Sets the background color to black.
        - text-white: Sets the default text color to white for contrast.
        - flex flex-col items-center justify-center: Uses flexbox to center the content
          both horizontally and vertically within the viewport.
        - p-4: Adds padding on all sides, crucial for responsiveness on smaller screens.
      */}
      <div className="min-h-screen bg-black text-white flex flex-col items-center justify-center p-4">
        {/*
          This inner div acts as a content wrapper.
          - w-full: Ensures it takes full width on small screens.
          - max-w-screen-lg: Limits the maximum width of the content on larger screens
            to improve readability and visual appeal, keeping it centered.
        */}
        <div className="w-full max-w-screen-lg">{children}</div>
      </div>
    </ThemeProviderContext.Provider>
  );
}

export const useTheme = () => {
  const context = useContext(ThemeProviderContext);

  if (context === undefined)
    throw new Error("useTheme must be used within a ThemeProvider");

  return context;
};
