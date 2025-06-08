interface Props {
  balance: number;
  changePercent: number;
}

export function DashboardBalance({ balance, changePercent }: Props) {
  return (
    <div className="text-center">
      <h2 className="text-4xl font-bold">{balance} PLN</h2>
      <p className="text-green-500">+{changePercent} (YTD)</p>
    </div>
  );
}
