interface Props {
  balance: number;
  changePercent: number;
}

export function DashboardBalance({ balance, changePercent }: Props) {
  const formattedBalance = new Intl.NumberFormat("pl-PL", {
    style: "currency",
    currency: "PLN",
    minimumFractionDigits: 2,
  }).format(balance);

  const formattedPercent = `${changePercent > 0 ? "+" : ""}${changePercent.toFixed(2)}%`;

  const percentColor =
    changePercent > 0
      ? "text-green-500"
      : changePercent < 0
        ? "text-red-500"
        : "text-zinc-400";

  return (
    <div className="text-center">
      <h2 className="text-4xl font-bold">{formattedBalance}</h2>
      <p className={percentColor}>{formattedPercent} (YTD)</p>
    </div>
  );
}
