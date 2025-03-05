import React, { useEffect, useState } from "react";
import {
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { Check } from "lucide-react";
import { Checkout, SubscriptionPlans } from "@/queries/payments";
import { useQuery, useMutation } from "@apollo/client";

type SubscriptionDialogProps = {
  showOpen?: boolean;
  onClose: () => void;
  onSubscribe: (id: string) => void;
};

export default function SubscriptionDialog({
  showOpen,
  onClose,
}: SubscriptionDialogProps) {
  const [selectedPlan, setSelectedPlan] = useState<string>("");

  const [checkoutMutation] = useMutation(Checkout);
  const { data, loading } = useQuery(SubscriptionPlans, {
    skip: !showOpen,
  });

  useEffect(() => {
    if (
      !loading &&
      data &&
      data.subscriptionPlans?.length > 0 &&
      selectedPlan === ""
    ) {
      setSelectedPlan(data.subscriptionPlans[0].id);
    }
  }, [data, loading, selectedPlan]);

  const handleSubscribe = async (e: React.SyntheticEvent) => {
    e.preventDefault();
    if (!selectedPlan) {
      return;
    }

    const { data } = await checkoutMutation({
      variables: {
        id: selectedPlan,
      },
    });

    if (data?.checkoutSubscriptionPlan) {
      window.location.href = data.checkoutSubscriptionPlan.url;
    }

    onClose();
  };

  return (
    <DialogContent className="sm:max-w-[550px]">
      <DialogHeader>
        <DialogTitle>Choose your subscription plan</DialogTitle>
        <DialogDescription>
          Select the plan that works best for you. You can always change this
          later.
        </DialogDescription>
      </DialogHeader>
      <form onSubmit={handleSubscribe}>
        <input type="hidden" name="plan" value={selectedPlan} />

        <Card>
          <CardHeader>
            <CardTitle>Subscription Plans</CardTitle>
            <CardDescription>
              We offer flexible options to suit your needs
            </CardDescription>
          </CardHeader>
          <CardContent className="grid gap-6">
            <div className="grid grid-cols-2 gap-4">
              {data?.subscriptionPlans.map((plan) => (
                <PlanCard
                  key={plan.id}
                  title={plan.name}
                  price={`$${plan.priceCents / 100}`}
                  description={`per ${plan.interval}`}
                  isSelected={selectedPlan === plan.id}
                  onClick={() => setSelectedPlan(plan.id)}
                />
              ))}
            </div>
          </CardContent>
          <CardFooter>
            <Button type="submit" className="w-full">
              Subscribe Now
            </Button>
          </CardFooter>
        </Card>
      </form>
      <DialogFooter className="sm:justify-start">
        <DialogDescription>
          Prices are in USD. By subscribing, you agree to our Terms of Service.
        </DialogDescription>
      </DialogFooter>
    </DialogContent>
  );
}

interface PlanCardProps {
  title: string;
  price: string;
  description: string;
  isSelected: boolean;
  onClick: () => void;
}

function PlanCard({
  title,
  price,
  description,
  isSelected,
  onClick,
}: PlanCardProps) {
  return (
    <Card
      className={cn(
        "cursor-pointer transition-colors",
        isSelected
          ? "border-primary bg-primary text-primary-foreground"
          : "hover:border-primary",
      )}
      onClick={onClick}
    >
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="text-3xl font-bold">{price}</div>
        <p className="text-sm opacity-90">{description}</p>
      </CardContent>
      <CardFooter>
        {isSelected && (
          <div className="rounded-full bg-background p-1">
            <Check className="h-4 w-4 text-primary" />
          </div>
        )}
      </CardFooter>
    </Card>
  );
}
