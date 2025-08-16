import { fetchHouseholdUsers } from "@/app/api/fetchHouseholdUsers";
import { TLoginUser } from "@/app/api/fetchLoginUser";
import HouseholdClientPage from "@/app/household/components/householdClientPage";

export default async function HouseholdPage() {
  // Fetching data on the server component
  const householdUsers: TLoginUser[] = await fetchHouseholdUsers();

  return (
    <main>
      <h1>世帯管理</h1>
      <HouseholdClientPage initialUsers={householdUsers} />
    </main>
  );
}
