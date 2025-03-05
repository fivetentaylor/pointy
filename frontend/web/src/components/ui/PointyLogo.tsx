import Image from "next/image";

const PointyLogo = function () {
  return (
    <>
      {/* Light mode logo */}
      <Image
        className="dark:hidden block w-28 h-auto"
        src="/images/pointy.png"
        alt="Pointy Logo Light"
        width={112}
        height={26}
      />

      {/* Dark mode logo */}
      <Image
        className="dark:block hidden w-28 h-auto"
        src="/images/pointy.png"
        alt="Pointy Logo Dark"
        width={112}
        height={26}
      />
    </>
  );
};

export default PointyLogo;
