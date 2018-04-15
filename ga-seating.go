package main

import (
  "fmt"
  "math/rand"
  "time"
  "sort"
  //"strings"
)

type byScore []int

// array of tables, each table is defined with a number of seats
// tables can be in arbitrary order
var tables = []int{4,4,4,3,3}

// number of times to change seating
var rotations = 9
var source = rand.NewSource(time.Now().UnixNano())
var pop_size = 10000
var gen_size = 1000
var mutation_permille = 32
var scores []int
var member_count = 0

// ****************************************************************************
// Seatings
// ****************************************************************************
func generate_seating() []int {
  // simple seating where members sit in order
  seating := make( []int, member_count )
  for i := 0; i < member_count; i++ {
    seating[ i ] = i
  }
  return seating
}

func copy_seating( seat []int ) []int {
  seating := make( []int, member_count )
  for i := 0; i < member_count; i++ {
    seating[ i ] = seat[ i ]
  }
  return seating
}

func shuffle_seating( seating []int ){
  random := rand.New(source)
  // fisher-yates shuffle
  for i := len( seating ) - 1; i > 0; i-- {
     j := random.Intn(i + 1)
     seating[i], seating[j] = seating[j], seating[i]
  }
  order_tables( seating )
}

func count_meetings( seating []int, meetings []int ){
  // count the number of times each individual is put in the same group
  // as every other individual
  // this can be optimized by making 'meetings' a triangular array instead of a square one
  offset := 0
  for _, g := range tables {
    for i := 0; i < g; i++ {
      for j := i+1; j < g; j++ {
        if seating[i+offset] < seating[j+offset] {
          meetings[ member_count*seating[i+offset] + seating[j+offset] ]+=1
        }else{
          // because tables are sorted internally,
          // the member at i should always be smaller than the one at j
          // so this should never happen
          meetings[ member_count*seating[j+offset] + seating[i+offset] ]+=1
          fmt.Println( "Fail.  This should never happen.  See source.")
        }
      }
    }
    offset += g
  }
}

func print_seating( seating []int ) {
  fmt.Println( fmt.Sprint(seating), " ", evaluate_seating(seating) )
}

func order_tables( seating []int ){
  // sort each table in the seating
  offset := 0
  for _, g := range tables {
    for i := 0; i < g; i++ {
      sort.Ints( seating[ offset : (g + offset) ] )
    }
    offset += g
  }
}

func evaluate_seating( seating []int) int {
  // witness the fitness of a single seating
  fitness := 0

  // count the number of times each member is paired with another
  meetings := make( []int, member_count*member_count )
  count_meetings( seating, meetings )

  // score the seating based on the meeting count
  // todo: make a better heuristic for scoring
  for _,i := range meetings {
    if( i == 1 ){
      fitness += 1
    }
  }
  return fitness
}


// ****************************************************************************
// Schedules
// ****************************************************************************

func evaluate_schedule( sched [][]int ) int {
  // find a fitness for a schedule, favoring counts of one and lower counts
  fitness := 0
  meetings := make( []int, member_count*member_count )

  // optimize 0, it should always be the same
  count_meetings( sched[0], meetings )
  for i := 1; i < rotations; i++ {
    count_meetings( sched[ i ], meetings )
  }
  for _,i := range meetings {
    if( i == 1 ){
      fitness += 100
    }else if ( i > 1 ){
      fitness = 100 - i
    }
  }
  return fitness
}

func evaluate_schedule_simple( sched [][]int ) int {
  fitness := 0
  meetings := make( []int, member_count*member_count )

  // optimize 0, it should always be the same
  count_meetings( sched[0], meetings )
  for i := 1; i < rotations; i++ {
    count_meetings( sched[ i ], meetings )
  }
  for _,i := range meetings {
    if( i > 0 ){
      fitness += 1
    }
  }
  return fitness
}

func show_scores( sched [][]int ) {
  meetings := make( []int, member_count*member_count )

  count_meetings( sched[0], meetings )
  for i := 1; i < rotations; i++ {
    count_meetings( sched[ i ], meetings )
  }
  for i := 0; i < member_count; i++ {
    sum := 0
    for j := 0; j < member_count; j++ {
      sum += meetings[ (i * member_count ) + j ]
    }
    fmt.Println( meetings[ i*member_count : (i*(member_count)+member_count)] )
  }
}

func print_schedule( sched [][]int) {
  for i := 1; i < rotations; i++ {
    fmt.Println( fmt.Sprint(sched[i]) )
  }
  fmt.Println( "" )
}

func generate_schedule() [][]int {
  // a schedule is a list of several seatings
  sched := make( [][]int, rotations )
  sched[ 0 ] = generate_seating()
  for i := 1; i < rotations; i++ {
    sched[ i ] = generate_seating()
    shuffle_seating( sched[i] )
  }
  return sched
}

func copy_schedule( old [][]int ) [][]int {
  new := make( [][]int, rotations )
  for i := 0; i < rotations; i++ {
    newseat := make( []int, member_count )
    for j := 0; j < member_count; j++ {
      newseat[j] = old[i][j]
    }
    new[ i ] = newseat
  }
  return new
}

func mutate_schedule2( sched [][]int ){
  random := rand.New(source)
  rot := 1 + random.Intn( rotations - 1 )
  shuffle_seating( sched[ rot ] )
}


func mutate_schedule1( sched [][]int ){
  random := rand.New(source)
  // choose a seating (not the first!) to mutate
  rot := 1 + random.Intn( rotations - 1 )

  // choose two tables in the seating to switch between
  table1 := random.Intn( len (tables) )
  table2 := random.Intn( len (tables) )
  for ; table1 == table2; {
    table2 = random.Intn( len (tables) )
  }

  seat1 := 0
  seat2 := 0
  for i := 0; i < table1; i+= 1 {
    seat1 += tables[i]
  }
  for i := 0; i < table2; i+= 1 {
    seat2 += tables[i]
  }
  seat1 += random.Intn( tables[ table1 ] )
  seat2 += random.Intn( tables[ table2 ] )

  sched[ rot ][ seat1 ], sched[ rot ][ seat2 ] = sched[ rot ][ seat2 ], sched[ rot ][ seat1 ]

  order_tables( sched[ rot ] )
}

func stepped_hill_climb(){
  // boo lame
  score1 := 0
  score2 := 0
  best := 0
  s1 := generate_schedule()
  s2 := copy_schedule( s1 )

  for eval_depth := 1; eval_depth <= rotations; eval_depth++ {
    for i := 0; i < 400000; i++ {
      score1 = evaluate_schedule( s1 )
      score2 = evaluate_schedule( s2 )
      if( score1 > score2 ){
        s2 = copy_schedule( s1 )
        mutate_schedule1( s2 )
      }else{
        s1 = copy_schedule( s2 )
        mutate_schedule2( s1 )
      }
      if( score1 > best ){
        best = score1
        fmt.Println( "new best: ", best )
        print_schedule( s1 )
      }
      if( score2 > best ){
        best = score2
        fmt.Println( "new best: ", best )
        print_schedule( s2 )
      }
    }
  }
}

// needed for custom sort
func (s byScore) Len() int {
  return len( s )
}

// needed for custom sort
func (s byScore) Swap(i, j int) {
  s[i], s[j] = s[j], s[i]
}

// needed for custom sort
func (s byScore) Less(i, j int) bool {
  return scores[s[i]] < scores[s[j]]
}

func hotsteamylove( pu1, pu2 [][]int ) [][]int{
  // create a new individual from two other individuals
  random := rand.New(source)

  junior := copy_schedule( pu1 )

  // choose a seating to take from other parent 
  seating_i := 1 + random.Intn( rotations - 1 )
  junior[ seating_i ] = copy_seating( pu2[ seating_i ] )

  return junior
}

// ****************************************************************************
// ****************************************************************************
func ga(){
  // create a initial population of random schedules
  // each generation, calculate each individual's fitness
  // sort them by fitness, and replace the lowest scoring
  // individuals with (sometimes mutated) offspring
  random := rand.New(source)

  // the general population
  pop := make( [][][]int, pop_size)

  // the list of scores:
  //   scores[ i ] = score( pop[ i ] )
  scores = make ([]int, pop_size)

  // the population index ordered by score:
  //   score( pop[ ordered_scores[ i ] ] ) <= score( pop[ ordered_scores[ i + 1 ] ] )
  ordered_scores := make ([]int, pop_size)

  for i := 0; i < pop_size; i++ {
    pop[ i ] = generate_schedule()
  }
  best := 0
  goal := (member_count * (member_count -1 )) / 2
  sum := 0
  for generation := 0; best < goal; generation++{
    sum = 0
    for i := 0; i < pop_size; i++ {
      ordered_scores[ i ] = i
      //scores[ i ] = evaluate_schedule( pop[ i ] )
      scores[ i ] = evaluate_schedule_simple( pop[ i ] )
      sum += scores[ i ]
      if( scores[ i ] > best ){
        best = scores[ i ]
        fmt.Println( "" )
        fmt.Println( "Human/Goal/Score/Generation", evaluate_schedule_simple( pop[ i ] ), goal, best, generation )
        print_schedule( pop[ i ] )
        show_scores( pop[ i ] )
      }
    }
    if generation % 500 == 0 {
      fmt.Println( "Generation/Goal/Best/Average", generation, goal, best, sum / pop_size )
    }
    sort.Sort( byScore( ordered_scores ) )
    for i := 0; i < gen_size; i++ {
      // replace weakest members with offspring of randomly selected members
      pu1 := random.Intn( pop_size )
      pu2 := random.Intn( pop_size )
      pop[ ordered_scores[ i ] ] = hotsteamylove( pop[ pu1 ], pop[ pu2 ] )
      if random.Intn( 1000 ) < mutation_permille {
        mutate_schedule1( pop[ i ] )
      }else if random.Intn( 1000 ) < mutation_permille {
        mutate_schedule2( pop[ i ] )
      }
    }
  }
}

func main() {
  for _,g := range tables{
    member_count += g
  }
  //stepped_hill_climb()
  ga()
}

