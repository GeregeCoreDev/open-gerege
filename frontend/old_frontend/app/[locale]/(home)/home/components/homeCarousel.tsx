import { useState, useEffect } from 'react'
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
  type CarouselApi,
} from '@/components/ui/carousel'
import Autoplay from 'embla-carousel-autoplay'
import Image from 'next/image'

const carouselSlides = [
  'https://picsum.photos/1920/400?random=1',
  'https://picsum.photos/1920/400?random=2',
  'https://picsum.photos/1920/400?random=3',
  'https://picsum.photos/1920/400?random=4',
  'https://picsum.photos/1920/400?random=5',
  'https://picsum.photos/1920/400?random=6',
]

export function HomeCarousel() {
  const [api, setApi] = useState<CarouselApi>()
  const [current, setCurrent] = useState(0)

  useEffect(() => {
    if (!api) return

    setCurrent(api.selectedScrollSnap())

    api.on('select', () => {
      setCurrent(api.selectedScrollSnap())
    })
  }, [api])

  return (
    <div className="relative w-full">
      <Carousel
        setApi={setApi}
        opts={{
          align: 'start',
          loop: true,
        }}
        plugins={[
          Autoplay({
            delay: 5000,
          }),
        ]}
        className="w-full"
      >
        <CarouselContent>
          {carouselSlides.map((slide, index) => {
            return (
              <CarouselItem key={slide}>
                <div className="relative h-[300px] lg:h-[400px]">
                  <Image
                    src={slide}
                    alt={slide}
                    fill
                    priority={index === 0}
                    className="absolute inset-0 h-full w-full object-cover"
                  />
                </div>
              </CarouselItem>
            )
          })}
        </CarouselContent>

        <div className="hidden lg:block">
          <CarouselPrevious className="left-4 shadow-none" />
          <CarouselNext className="right-4 shadow-none" />
        </div>
      </Carousel>

      {/* Carousel Indicators */}
      <div className="absolute bottom-4 flex w-full justify-center gap-2">
        {carouselSlides.map((_, index) => (
          <button
            key={index}
            onClick={() => api?.scrollTo(index)}
            className={`h-2 rounded-full transition-all duration-300 ${current === index ? 'bg-primary w-8' : 'w-2 bg-black'
              }`}
            aria-label={`Go to slide ${index + 1}`}
          />
        ))}
      </div>
    </div>
  )
}
