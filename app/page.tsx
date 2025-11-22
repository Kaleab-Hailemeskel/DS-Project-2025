import { Link } from "@heroui/link";
import { Snippet } from "@heroui/snippet";
import { Code } from "@heroui/code";
import { button as buttonStyles } from "@heroui/theme";

import { siteConfig } from "@/config/site";
import { title, subtitle } from "@/components/primitives";
import { GithubIcon } from "@/components/icons";
import BlurredCard from "@/components/BlurredCard";
import MusicPlayer from "@/components/MusicPlayer"

export default function Home() {
  return (
    <section className="flex flex-col items-center justify-center gap-4 py-8 md:py-10">
      <div className="inline-block max-w-xl text-center justify-center">
        <span className={title()}>Explore&nbsp;</span>
        <span className={title({ color: "pink" })}>Music&nbsp;</span>
        <br />
        <span className={title()}>
          like never before
        </span>
        <div className={subtitle({ class: "mt-4" })}>
          Pop, Rap and everything in between
        </div>
      </div>

      <div className="flex gap-3">
        <Link
          isExternal
          className={buttonStyles({
            color: "danger",
            radius: "full",
            variant: "shadow",
          })}
          href={siteConfig.links.docs}
        >
          Start Listening Now
        </Link>
        
      </div>

      <div className="mt-8">
        <Snippet hideCopyButton hideSymbol variant="bordered">
          <span>
            Get started by editing <Code color="primary">app/page.tsx</Code>
          </span>
        </Snippet>
      </div>
      <BlurredCard></BlurredCard>
       <MusicPlayer></MusicPlayer>
    </section>
  );
}
