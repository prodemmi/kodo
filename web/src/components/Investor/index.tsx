import { ReactNode } from "react";
import { useAppState } from "../../states/app.state";

type GuardProps = {
  children: ReactNode;
};

function Investor({ children }: GuardProps) {
  const isInvestor = useAppState((s) => s.investor);
  return isInvestor ? <>{children}</> : null;
}

function Consumer({ children }: GuardProps) {
  const isInvestor = useAppState((s) => s.investor);
  return !isInvestor ? <>{children}</> : null;
}

export const RoleGuard = {
  Investor,
  Consumer,
};
