## Mecachis

![Ay mecachis!](./art/mecachis.png)

- Is your system running out of resources?
- Are your services unable to scale?
- Do you think that a good and clean architecture seems like too much?
- Is your API slow as hell because it does more than it should?

If you have answered all those questions with a "yes", fear no longer. This
brand-new package aims to someday provide those extra layers of 
_cachabstraction_ and indirection that your application needs. Remember: You
can install a cache in the frontend, in front of a web server or in between
any two services that need to communicate with each other. That way, no
one would ever know where the source of truth is but yeah, the app
runs smoothly.

Now, jokes apart, this is just another academic repo to learn. This
time I'm into caches and the wide range of strategies implemented out
there to adapt caching to a wider range of scenarios. What I would
like to achieve with this repo is gaining deeper knowledge on data 
structures and algorithms. Beyond that, it would be even nicer if:

- More than 5 strategies are implemented [3/5]
    - [x] LRU
    - [x] LFU
    - [x] LFRU
- Caches are distributed over the network
- Any kind of background persistence is achieved

## Credits

### Cache eviction algorithms

Buf of course, algorithms implemented here are not of my invention, but slightly modified
to my particular needs and influenced by my knowledge state, though. I was heavily inspired
by [this wikipedia entry](https://en.wikipedia.org/wiki/Cache_replacement_policies) about cache
eviction algos.

**External links and references**

- LFRU: M. Bilal and S. -G. Kang, "A Cache Management Scheme for Efficient Content Eviction and Replication in Cache Networks," in IEEE Access, vol. 5, pp. 1692-1701, 2017, doi: 10.1109/ACCESS.2017.2669344. **Paper**: https://arxiv.org/ftp/arxiv/papers/1702/1702.04078.pdf **Patent**: https://patentimages.storage.googleapis.com/60/c5/34/c94ab8b27e2f9d/US10819823.pdf